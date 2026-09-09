// Package mixrank is a streaming Go client for MixRank. All API surfaces share
// its credential redaction, cancellation, retry and validation concurrency rules.
package mixrank

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/flock"
	"github.com/nkulavic/mixrank/catalog"
)

const DefaultBaseURL = "https://api.mixrank.com/v2/json"
const DefaultMaxResponseBytes int64 = 16 << 20

type Client struct {
	key      string
	base     *url.URL
	http     *http.Client
	retries  int
	lockDir  string
	volatile bool
}
type Options struct {
	InMemoryCoordination bool // For stateless HTTP hosts; separate processes/hosts require caller coordination.
	BaseURL              string
	HTTPClient           *http.Client
	Retries              int
	LockDir              string
}

func New(apiKey string, opts Options) (*Client, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("MixRank credential missing; run mixrank auth login")
	}
	if strings.ContainsAny(apiKey, "\r\n/\\?#%") {
		return nil, errors.New("invalid API credential format")
	}
	base := opts.BaseURL
	if base == "" {
		base = DefaultBaseURL
	}
	u, err := url.Parse(base)
	if err != nil || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return nil, errors.New("invalid API base URL")
	}
	if u.Scheme != "https" && !(u.Scheme == "http" && (u.Hostname() == "127.0.0.1" || u.Hostname() == "localhost" || u.Hostname() == "::1")) {
		return nil, errors.New("API requires HTTPS (except loopback tests)")
	}
	h := opts.HTTPClient
	if h == nil {
		h = &http.Client{Timeout: 150 * time.Second}
	}
	copyClient := *h
	// Never forward the API key through redirects.
	copyClient.CheckRedirect = func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }
	dir := opts.LockDir
	if dir == "" {
		home, e := os.UserHomeDir()
		if e != nil {
			return nil, errors.New("cannot locate credential lock directory")
		}
		dir = filepath.Join(home, ".mixrank", "locks")
	}
	n := opts.Retries
	if n < 0 || n > 5 {
		return nil, errors.New("retries must be 0..5")
	}
	return &Client{key: apiKey, base: u, http: &copyClient, retries: n, lockDir: dir, volatile: opts.InMemoryCoordination}, nil
}

type Request struct {
	Parameters map[string]string `json:"parameters,omitempty"`
	Query      url.Values        `json:"query,omitempty"`
	Body       json.RawMessage   `json:"body,omitempty"`
	Upload     io.Reader         `json:"-"`
	Filename   string            `json:"filename,omitempty"`
	Refresh    bool              `json:"refresh,omitempty"`
}
type Response struct {
	StatusCode int           `json:"status"`
	Header     http.Header   `json:"headers"`
	Body       io.ReadCloser `json:"-"`
}

func (r *Response) JSON(limit int64) (any, error) {
	if limit <= 0 {
		limit = DefaultMaxResponseBytes
	}
	b, e := io.ReadAll(io.LimitReader(r.Body, limit+1))
	if e != nil {
		return nil, errors.New("response read interrupted")
	}
	if int64(len(b)) > limit {
		return nil, errors.New("response exceeds limit; stream to a file or narrow query")
	}
	if len(bytes.TrimSpace(b)) == 0 {
		return nil, nil
	}
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	var v any
	if e = d.Decode(&v); e != nil {
		return nil, errors.New("response is not JSON; use raw output")
	}
	return v, nil
}

type APIError struct {
	Status       int
	Method, Path string
	Body         json.RawMessage
	Uncertain    bool
}

func (e *APIError) Error() string {
	if e.Uncertain {
		return fmt.Sprintf("MixRank %s %s: HTTP %d; outcome may be uncertain, inspect existing job/result before retry", e.Method, e.Path, e.Status)
	}
	return fmt.Sprintf("MixRank %s %s: HTTP %d", e.Method, e.Path, e.Status)
}

var keyPath = regexp.MustCompile(`(/v2/json/)[^/\s?"<>]+`)

func (c *Client) Redact(s string) string {
	s = strings.ReplaceAll(s, c.key, "[REDACTED]")
	s = strings.ReplaceAll(s, url.PathEscape(c.key), "[REDACTED]")
	return keyPath.ReplaceAllString(s, "${1}[REDACTED]")
}
func (c *Client) Call(ctx context.Context, id string, in Request) (*Response, error) {
	op, e := catalog.Lookup(id)
	if e != nil {
		return nil, e
	}
	if op.Encoding == "form" && len(in.Body) > 0 {
		var fields map[string]json.RawMessage
		if e := json.Unmarshal(in.Body, &fields); e != nil {
			return nil, errors.New("form body must be an object")
		}
		params := map[string]string{}
		for k, v := range fields {
			var str string
			if json.Unmarshal(v, &str) == nil {
				params[k] = str
			} else {
				params[k] = string(v)
			}
		}
		for k, v := range in.Parameters {
			params[k] = v
		}
		in.Parameters = params
		in.Body = nil
	}
	path := op.Path
	q := url.Values{}
	for k, v := range in.Query {
		q[k] = append([]string(nil), v...)
	}
	form := url.Values{}
	for _, p := range op.Parameters {
		v, ok := in.Parameters[p.Name]
		if !ok {
			v = q.Get(p.Name)
		}
		if p.In == "file" {
			if p.Required && in.Upload == nil {
				return nil, errors.New("operation requires a file upload")
			}
			continue
		}
		if p.Required && v == "" {
			return nil, fmt.Errorf("missing parameter %s", p.Name)
		}
		if p.In == "path" {
			if v == "" || v == "." || v == ".." {
				return nil, fmt.Errorf("invalid path parameter %s", p.Name)
			}
			path = strings.ReplaceAll(path, "{"+p.Name+"}", url.PathEscape(v))
			q.Del(p.Name)
		} else if ok {
			if p.In == "form" {
				form.Set(p.Name, v)
			} else {
				q.Set(p.Name, v)
			}
		}
	}
	for k, v := range in.Parameters {
		known := false
		for _, p := range op.Parameters {
			if p.Name == k {
				known = true
				break
			}
		}
		if !known {
			if op.Encoding == "form" {
				form.Set(k, v)
			} else {
				q.Set(k, v)
			}
		}
	}
	var body io.Reader
	contentType := ""
	var pipe *io.PipeReader
	switch op.Encoding {
	case "json":
		if len(in.Body) == 0 {
			in.Body = json.RawMessage(`{}`)
		}
		if !json.Valid(in.Body) {
			return nil, errors.New("body must be valid JSON")
		}
		body = bytes.NewReader(in.Body)
		contentType = "application/json"
	case "form":
		for k, v := range q {
			form[k] = v
		}
		q = url.Values{}
		if len(in.Body) > 0 {
			var obj map[string]json.RawMessage
			if e := json.Unmarshal(in.Body, &obj); e != nil {
				return nil, errors.New("form body must be a JSON object")
			}
			for k, v := range obj {
				var s string
				if json.Unmarshal(v, &s) == nil {
					form.Set(k, s)
				} else {
					form.Set(k, string(v))
				}
			}
		}
		body = strings.NewReader(form.Encode())
		contentType = "application/x-www-form-urlencoded"
	case "multipart":
		r, w := io.Pipe()
		pipe = r
		mw := multipart.NewWriter(w)
		contentType = mw.FormDataContentType()
		body = r
		name := in.Filename
		if name == "" {
			name = "input.txt"
		}
		go func() {
			part, e := mw.CreateFormFile("file", filepath.Base(name))
			if e == nil {
				_, e = io.Copy(part, in.Upload)
			}
			if e == nil {
				e = mw.Close()
			}
			_ = w.CloseWithError(e)
		}()
	}
	if pipe != nil {
		defer pipe.Close()
	}
	return c.do(ctx, op.Method, path, q, body, contentType, in.Refresh, &op)
}

// Raw still enforces validation serialization and refresh/retry semantics.
func (c *Client) Raw(ctx context.Context, method, path string, q url.Values, body io.Reader, contentType string, refresh bool) (*Response, error) {
	method = strings.ToUpper(method)
	if method != "GET" && method != "POST" && method != "DELETE" && method != "PUT" && method != "PATCH" {
		return nil, errors.New("unsupported HTTP method")
	}
	if !strings.HasPrefix(path, "/") || strings.ContainsAny(path, "%?#\\") || strings.Contains(path, "//") {
		return nil, errors.New("path must be an unescaped relative API path without query or fragment")
	}
	for _, p := range strings.Split(path, "/") {
		if p == "." || p == ".." {
			return nil, errors.New("path traversal is not allowed")
		}
	}
	path = strings.TrimSuffix(path, "/")
	op, ok := catalog.Match(method, path)
	if !ok {
		return c.do(ctx, method, path, q, body, contentType, refresh, nil)
	}
	return c.do(ctx, method, path, q, body, contentType, refresh, &op)
}

var localLocks sync.Map

func (c *Client) validationLock(ctx context.Context) (func(), error) {
	sum := fmt.Sprintf("%x", sha256.Sum256([]byte(c.key)))
	v, _ := localLocks.LoadOrStore(sum, make(chan struct{}, 1))
	ch := v.(chan struct{})
	select {
	case ch <- struct{}{}:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	release := func() { <-ch }
	if c.volatile {
		return release, nil
	}
	if e := os.MkdirAll(c.lockDir, 0700); e != nil {
		release()
		return nil, errors.New("cannot create validation lock directory")
	}
	f := flock.New(filepath.Join(c.lockDir, sum+".lock"))
	ok, e := f.TryLockContext(ctx, 50*time.Millisecond)
	if e != nil || !ok {
		release()
		return nil, errors.New("cannot acquire validation lock")
	}
	return func() { _ = f.Unlock(); release() }, nil
}

type lockedBody struct {
	io.ReadCloser
	once    sync.Once
	release func()
}

func (b *lockedBody) Close() error { e := b.ReadCloser.Close(); b.once.Do(b.release); return e }
func (c *Client) do(ctx context.Context, method, path string, q url.Values, body io.Reader, contentType string, refresh bool, op *catalog.Operation) (*Response, error) {
	if path == "/email/validate" && (method != "GET" || body != nil) {
		return nil, errors.New("single email validation requires GET query parameters")
	}
	query := url.Values{}
	for k, v := range q {
		query[k] = append([]string(nil), v...)
	}
	for _, k := range []string{"strategy", "webhook", "private"} {
		if len(query[k]) > 1 {
			return nil, fmt.Errorf("duplicate %s parameters are not supported", k)
		}
	}
	safe := op != nil && op.SafeRetry
	if op != nil {
		for _, p := range op.Parameters {
			if p.Name == "strategy" {
				if query.Get("strategy") == "" {
					if strings.HasPrefix(path, "/linkedin/") {
						query.Set("strategy", "cached")
					}
				}
				s := query.Get("strategy")
				if s == "" {
					s = p.Default
				}
				if strings.HasPrefix(path, "/linkedin/") && s != "cached" && !refresh {
					return nil, errors.New("strategy may fetch live data; pass explicit refresh")
				}
				if s != "cached" {
					safe = false
				}
			}
		}
	}
	if query.Get("strategy") != "" && query.Get("strategy") != "cached" {
		safe = false
	}

	if path == "/email/validate/bulk-job" && method == "POST" && (query.Get("private") == "t" || query.Get("private") == "true" || query.Get("private") == "1") && query.Get("strategy") != "fetch" {
		return nil, errors.New("private bulk jobs require strategy=fetch")
	}
	var unlock func()
	if path == "/email/validate" {
		var err error
		unlock, err = c.validationLock(ctx)
		if err != nil {
			return nil, err
		}
		if e := c.validatePending(); e != nil {
			unlock()
			return nil, e
		}
		defer func() {
			if unlock != nil {
				unlock()
			}
		}()
	}
	u := *c.base
	decodedPath, err := url.PathUnescape(path)
	if err != nil {
		return nil, errors.New("invalid path encoding")
	}
	u.Path = strings.TrimRight(c.base.Path, "/") + "/" + c.key + decodedPath
	u.RawPath = strings.TrimRight(c.base.EscapedPath(), "/") + "/" + url.PathEscape(c.key) + path
	u.RawQuery = query.Encode()
	req, e := http.NewRequestWithContext(ctx, method, u.String(), body)
	if e != nil {
		return nil, errors.New("cannot construct request")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "mixrank-go/0.2.0")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	attempts := 1
	if safe && (body == nil || req.GetBody != nil) {
		attempts += c.retries
	}
	for attempt := 0; attempt < attempts; attempt++ {
		if attempt > 0 {
			if req.GetBody != nil {
				req.Body, _ = req.GetBody()
			}
		}
		res, err := c.http.Do(req)
		if err != nil {
			if path == "/email/validate" {
				_ = c.setPending(ValidationPending{AcceptedAt: time.Now().UTC(), Uncertain: true})
			}
			return nil, &APIError{Status: 0, Method: method, Path: c.Redact(path), Uncertain: !safe}
		}
		if res.StatusCode >= 200 && res.StatusCode < 300 {
			rc := res.Body
			if path == "/email/validate" && res.StatusCode == 202 {
				b, readErr := io.ReadAll(io.LimitReader(res.Body, 65537))
				res.Body.Close()
				var accepted struct {
					ID string `json:"id"`
				}
				uncertain := readErr != nil || len(b) > 65536 || json.Unmarshal(b, &accepted) != nil || accepted.ID == ""
				if e := c.setPending(ValidationPending{ID: accepted.ID, AcceptedAt: time.Now().UTC(), Uncertain: uncertain}); e != nil {
					return nil, e
				}
				rc = io.NopCloser(bytes.NewReader(b))
			}
			if unlock != nil {
				rc = &lockedBody{ReadCloser: rc, release: unlock}
				unlock = nil
			}
			return &Response{StatusCode: res.StatusCode, Header: res.Header.Clone(), Body: rc}, nil
		}
		b, _ := io.ReadAll(io.LimitReader(res.Body, 64<<10))
		res.Body.Close()
		retry := res.StatusCode == 429 || res.StatusCode == 502 || res.StatusCode == 503 || res.StatusCode == 504
		if retry && attempt+1 < attempts {
			delay := time.Duration(1<<attempt) * time.Second
			if n, e := strconv.Atoi(res.Header.Get("Retry-After")); e == nil && n > 0 {
				delay = time.Duration(n) * time.Second
			} else if t, e := http.ParseTime(res.Header.Get("Retry-After")); e == nil {
				delay = time.Until(t)
			}
			if delay < 0 {
				delay = 0
			}
			if delay > 60*time.Second {
				return nil, &APIError{Status: res.StatusCode, Method: method, Path: path, Body: json.RawMessage(c.Redact(string(b)))}
			}
			timer := time.NewTimer(delay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			case <-timer.C:
			}
			continue
		}
		return nil, &APIError{Status: res.StatusCode, Method: method, Path: c.Redact(path), Body: json.RawMessage(c.Redact(string(b))), Uncertain: !safe && (res.StatusCode >= 500 || res.StatusCode == 408)}
	}
	return nil, errors.New("request failed")
}
