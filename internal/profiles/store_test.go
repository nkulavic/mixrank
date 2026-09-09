package profiles

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestReviewRefreshAndAttribution(t *testing.T) {
	s := Store{Root: t.TempDir()}
	src := Source{Kind: "website", Location: "https://example.test"}
	task, e := s.Draft("product", []Source{src})
	if e != nil {
		t.Fatal(e)
	}
	if _, e = os.Stat(task); e != nil {
		t.Fatal(e)
	}
	if e = s.Accept("product"); e == nil {
		t.Fatal("accepted empty analysis")
	}
	p, _ := s.Show("product", true)
	p.Positioning = "Synthetic backup tool"
	p.Evidence = []Evidence{{Statement: "Supports synthetic backups", Kind: "implementation_evidence", Source: src.Location, ObservedAt: time.Now().UTC().Format(time.RFC3339)}}
	d, _ := s.Dir("product")
	if e = WritePrivate(filepath.Join(d, "draft.json"), p); e != nil {
		t.Fatal(e)
	}
	if e = s.Accept("product"); e != nil {
		t.Fatal(e)
	}
	active, _ := s.Show("product", false)
	if e = s.Select("product", ""); e != nil {
		t.Fatal(e)
	}
	if name, e := s.Selected("", ""); e != nil || name != "product" {
		t.Fatal(name, e)
	}
	if _, e = s.Draft("product", nil); e != nil {
		t.Fatal(e)
	}
	p, _ = s.Show("product", true)
	p.Positioning = "New positioning"
	p.Evidence[0].Source = "https://unselected.test"
	WritePrivate(filepath.Join(d, "draft.json"), p)
	if s.Accept("product") == nil {
		t.Fatal("accepted unattributed evidence")
	}
	current, _ := s.Show("product", false)
	if current.Version != active.Version {
		t.Fatal("refresh replaced active before acceptance")
	}
	if _, e = s.Dir("../../escape"); e == nil {
		t.Fatal("path traversal")
	}
	if runtime.GOOS != "windows" {
		if info, e := os.Stat(filepath.Join(d, "active.json")); e == nil && info.Mode().Perm()&0077 != 0 {
			t.Fatal("private permissions")
		}
	}
}
