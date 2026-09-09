package mixrank

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nkulavic/mixrank/catalog"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func testClient(t *testing.T, h http.HandlerFunc) *Client {
	t.Helper()
	s := httptest.NewServer(h)
	t.Cleanup(s.Close)
	c, e := New("synthetic-test-credential", Options{BaseURL: s.URL, LockDir: t.TempDir()})
	if e != nil {
		t.Fatal(e)
	}
	return c
}
func TestCatalogAndRawMatching(t *testing.T) {
	c := catalog.Read()
	if len(c.Operations) != 102 {
		t.Fatalf("operations %d", len(c.Operations))
	}
	seen := map[string]bool{}
	for _, op := range c.Operations {
		key := op.Method + op.Path
		if seen[key] {
			t.Fatal("duplicate", key)
		}
		seen[key] = true
		if len(op.Sources) == 0 || op.Verified == "" {
			t.Fatal("missing provenance", op.ID)
		}
		for _, p := range op.Parameters {
			if p.In == "path" && !p.Required {
				t.Fatal("optional path", p.Name)
			}
		}
	}
	op, ok := catalog.Match("GET", "/companies/match")
	if !ok || op.ID != "get_companies_match" || op.SafeRetry {
		t.Fatal("raw route lost semantic retry policy")
	}
}
func TestEncodingsNumbersAndCandidates(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/synthetic-test-credential/person/match":
			if r.URL.Query().Get("first_name") != "Alex" {
				t.Error("query encoding")
			}
			fmt.Fprint(w, `[{"id":9007199254740993,"score":0.5},{"id":2,"score":0.5}]`)
		case "/synthetic-test-credential/identity/ingest":
			if e := r.ParseForm(); e != nil {
				t.Error(e)
			}
			if r.Form.Get("source_key") != "one" || r.Form.Get("accounts") != `{"email_address":{"email":"example@example.test"}}` {
				t.Error("form body", r.Form)
			}
			fmt.Fprint(w, `{}`)
		case "/synthetic-test-credential/email/validate/bulk-job":
			if e := r.ParseMultipartForm(1 << 20); e != nil {
				t.Fatal(e)
			}
			f, _, e := r.FormFile("file")
			if e != nil {
				t.Fatal(e)
			}
			defer f.Close()
			b, _ := io.ReadAll(f)
			if string(b) != "a@example.test\nb@example.test\n" {
				t.Error("upload")
			}
			fmt.Fprint(w, `{"id":"job-one","completed_at":null}`)
		default:
			t.Error(r.URL.Path)
			http.Error(w, "bad", 400)
		}
	})
	r, e := c.GetPersonMatch(context.Background(), Request{Parameters: map[string]string{"first_name": "Alex", "page_size": "10"}})
	if e != nil {
		t.Fatal(e)
	}
	v, e := r.JSON(1024)
	r.Body.Close()
	if e != nil {
		t.Fatal(e)
	}
	a := v.([]any)
	if len(a) != 2 || a[0].(map[string]any)["id"].(json.Number).String() != "9007199254740993" {
		t.Fatal(v)
	}
	r, e = c.PostIdentityIngest(context.Background(), Request{Body: json.RawMessage(`{"source":"test","source_key":"one","accounts":{"email_address":{"email":"example@example.test"}}}`)})
	if e != nil {
		t.Fatal(e)
	}
	r.Body.Close()
	r, e = c.PostEmailValidateBulkJob(context.Background(), Request{Upload: strings.NewReader("a@example.test\nb@example.test\n")})
	if e != nil {
		t.Fatal(e)
	}
	r.Body.Close()
}
func TestRawValidationCannotBypass(t *testing.T) {
	var calls atomic.Int32
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) { calls.Add(1); fmt.Fprint(w, `{}`) })
	for _, p := range []string{"//email/validate", "/email/%76alidate", "/x/../email/validate", "/email/validate?x=y", "/email/validate\\"} {
		if _, e := c.Raw(context.Background(), "GET", p, nil, nil, "", false); e == nil {
			t.Fatal("unsafe raw path", p)
		}
	}
	for _, q := range []url.Values{{"strategy": {"cached", "fetch"}}, {"webhook": {"https://example.test/callback", "https://example.test/other"}}} {
		if _, e := c.Raw(context.Background(), "GET", "/email/validate", q, nil, "", false); e == nil {
			t.Fatal("unsafe parameters")
		}
	}
	if calls.Load() != 0 {
		t.Fatal("unexpected requests")
	}
}
func TestLiveRefreshSemanticRetriesAndRedaction(t *testing.T) {
	var calls atomic.Int32
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		http.Error(w, "https://api.mixrank.com/v2/json/synthetic-test-credential/linkedin/profile", 503)
	})
	c.retries = 2
	_, e := c.GetLinkedinProfile(context.Background(), Request{Parameters: map[string]string{"strategy": "fetch"}})
	if e == nil || calls.Load() != 0 {
		t.Fatal("refresh not gated")
	}
	_, e = c.GetLinkedinProfile(context.Background(), Request{Parameters: map[string]string{"strategy": "fetch"}, Refresh: true})
	if calls.Load() != 1 {
		t.Fatal("live fetch replayed")
	}
	var api *APIError
	if !errors.As(e, &api) || !api.Uncertain || strings.Contains(string(api.Body), c.key) || strings.Contains(e.Error(), c.key) {
		t.Fatal("redaction or uncertainty", e)
	}
}
func TestSingleValidationSerializedAcrossGoroutines(t *testing.T) {
	var active, max atomic.Int32
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		n := active.Add(1)
		for old := max.Load(); n > old; old = max.Load() {
			if max.CompareAndSwap(old, n) {
				break
			}
		}
		time.Sleep(15 * time.Millisecond)
		active.Add(-1)
		fmt.Fprint(w, `{"validity":"ambiguous"}`)
	})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, e := c.Raw(context.Background(), "GET", "/email/validate", url.Values{"email": {"example@example.test"}}, nil, "", false)
			if e != nil {
				t.Error(e)
				return
			}
			io.Copy(io.Discard, r.Body)
			r.Body.Close()
		}()
	}
	wg.Wait()
	if max.Load() != 1 {
		t.Fatalf("concurrent validations %d", max.Load())
	}
}
func TestProcessValidationHelper(t *testing.T) {
	if os.Getenv("MIXRANK_TEST_CHILD") != "1" {
		return
	}
	c, e := New("synthetic-test-credential", Options{BaseURL: os.Getenv("MIXRANK_TEST_URL"), LockDir: os.Getenv("MIXRANK_TEST_LOCK")})
	if e != nil {
		t.Fatal(e)
	}
	r, e := c.GetEmailValidate(context.Background(), Request{Parameters: map[string]string{"email": "example@example.test"}})
	if e != nil {
		t.Fatal(e)
	}
	io.Copy(io.Discard, r.Body)
	r.Body.Close()
}
func TestValidationAcrossProcesses(t *testing.T) {
	var active, max atomic.Int32
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := active.Add(1)
		for old := max.Load(); n > old; old = max.Load() {
			if max.CompareAndSwap(old, n) {
				break
			}
		}
		time.Sleep(80 * time.Millisecond)
		active.Add(-1)
		fmt.Fprint(w, `{}`)
	}))
	defer s.Close()
	dir := t.TempDir()
	cmds := []*exec.Cmd{}
	for i := 0; i < 3; i++ {
		cmd := exec.Command(os.Args[0], "-test.run=^TestProcessValidationHelper$")
		cmd.Env = append(os.Environ(), "MIXRANK_TEST_CHILD=1", "MIXRANK_TEST_URL="+s.URL, "MIXRANK_TEST_LOCK="+dir)
		if e := cmd.Start(); e != nil {
			t.Fatal(e)
		}
		cmds = append(cmds, cmd)
	}
	for _, cmd := range cmds {
		if e := cmd.Wait(); e != nil {
			t.Fatal(e)
		}
	}
	if max.Load() != 1 {
		t.Fatal("cross-process concurrency", max.Load())
	}
}
func TestOffsetPaginationAndPartialResults(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("offset") == "2" {
			http.Error(w, "license", 403)
			return
		}
		fmt.Fprint(w, `[{"id":1},{"id":2}]`)
	})
	pages := 0
	e := c.Pages(context.Background(), "get_companies", Request{Parameters: map[string]string{"page_size": "2"}}, 3, func(_ int, v any) error { pages++; return nil })
	if e == nil || pages != 1 {
		t.Fatal("partial pagination lost", e, pages)
	}
}
func TestWebhookVerification(t *testing.T) {
	now := time.Now()
	ts := fmt.Sprint(now.Unix())
	body := []byte(`{"id":"one","data":{}}`)
	mac := hmac.New(sha256.New, []byte("synthetic"))
	mac.Write([]byte("one." + ts + "."))
	mac.Write(body)
	h := http.Header{"Webhook-Id": []string{"one"}, "Webhook-Timestamp": []string{ts}, "Webhook-Signature": []string{"v1," + base64.StdEncoding.EncodeToString(mac.Sum(nil))}}
	if e := VerifyWebhook("synthetic", h, body, now, 5*time.Minute); e != nil {
		t.Fatal(e)
	}
	if VerifyWebhook("synthetic", h, append(body, ' '), now, 5*time.Minute) == nil {
		t.Fatal("accepted tamper")
	}
	if VerifyWebhook("synthetic", h, body, now.Add(-10*time.Minute), 5*time.Minute) == nil {
		t.Fatal("accepted future")
	}
}
func TestPathEscapingAndCancel(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.RequestURI, "hello%20world") {
			t.Error(r.RequestURI)
		}
		fmt.Fprint(w, `{}`)
	})
	r, e := c.PostTwitterQueryByQuery(context.Background(), Request{Parameters: map[string]string{"query": "hello world", "search_key": "one"}})
	if e != nil {
		t.Fatal(e)
	}
	r.Body.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, e = c.GetEcho(ctx, Request{})
	if e == nil {
		t.Fatal("cancellation ignored")
	}
}
func TestResponseBoundsAndEmpty(t *testing.T) {
	r := &Response{Body: io.NopCloser(strings.NewReader(`{"value":123456789}`))}
	if _, e := r.JSON(5); e == nil {
		t.Fatal("bound ignored")
	}
	r = &Response{StatusCode: 204, Body: io.NopCloser(bytes.NewReader(nil))}
	v, e := r.JSON(20)
	if e != nil || v != nil {
		t.Fatal("pending body", v, e)
	}
}
func TestSafeRetry(t *testing.T) {
	var count atomic.Int32
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if count.Add(1) == 1 {
			http.Error(w, "temporary", 503)
			return
		}
		fmt.Fprint(w, `{}`)
	})
	c.retries = 1
	r, e := c.GetEcho(context.Background(), Request{})
	if e != nil {
		t.Fatal(e)
	}
	r.Body.Close()
	if count.Load() != 2 {
		t.Fatal("safe retry")
	}
}
func TestMultipartFailureDoesNotLeakGoroutine(t *testing.T) {
	c, e := New("test", Options{BaseURL: "http://127.0.0.1:1", LockDir: filepath.Join(t.TempDir(), "locks")})
	if e != nil {
		t.Fatal(e)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, e = c.PostEmailValidateBulkJob(ctx, Request{Upload: strings.NewReader(strings.Repeat("x", 1<<20))})
	if e == nil {
		t.Fatal("expected failure")
	}
}

func TestAsyncValidationNeedsMatchingSignedCompletion(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(202)
		fmt.Fprint(w, `{"id":"async-one","status":"accepted"}`)
	})
	r, e := c.GetEmailValidate(context.Background(), Request{Parameters: map[string]string{"email": "a@example.test", "webhook": "https://callback.example/validation"}})
	if e != nil {
		t.Fatal(e)
	}
	r.Body.Close()
	if _, e = c.GetEmailValidate(context.Background(), Request{Parameters: map[string]string{"email": "b@example.test"}}); e == nil {
		t.Fatal("async work allowed parallel call")
	}
	now := time.Now()
	ts := fmt.Sprint(now.Unix())
	b := []byte(`{"id":"async-one","data":{"validity":"valid"}}`)
	mac := hmac.New(sha256.New, []byte(c.key))
	mac.Write([]byte("async-one." + ts + "."))
	mac.Write(b)
	h := http.Header{"Webhook-Id": []string{"async-one"}, "Webhook-Timestamp": []string{ts}, "Webhook-Signature": []string{"v1," + base64.StdEncoding.EncodeToString(mac.Sum(nil))}}
	if e = c.AcknowledgeValidationWebhook(context.Background(), h, b, now, 5*time.Minute); e != nil {
		t.Fatal(e)
	}
	if p, e := c.ValidationStatus(); e != nil || p != nil {
		t.Fatal("gate not cleared", p, e)
	}
}
