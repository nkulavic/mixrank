//go:build darwin

package credentials

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// Opt-in OS integration uses an isolated service, never the real credential.
func TestMacOSKeychainLifecycle(t *testing.T) {
	if os.Getenv("MIXRANK_TEST_KEYCHAIN") != "1" {
		t.Skip("opt-in OS credential test")
	}
	service := fmt.Sprintf("com.nkulavic.mixrank.test.%d", time.Now().UnixNano())
	defer exec.Command("/usr/bin/security", "delete-generic-password", "-a", "test", "-s", service).Run()
	for _, value := range []string{"synthetic-first", "synthetic-replacement"} {
		cmd := exec.Command("/usr/bin/security", "-i")
		cmd.Stdin = strings.NewReader("add-generic-password -U -a test -s " + service + " -w \"" + value + "\"\n")
		if _, e := cmd.CombinedOutput(); e != nil {
			t.Fatal("vault save failed")
		}
		b, e := exec.Command("/usr/bin/security", "find-generic-password", "-a", "test", "-s", service, "-w").Output()
		if e != nil || strings.TrimSpace(string(b)) != value {
			t.Fatal("vault read/replace failed")
		}
	}
	if e := exec.Command("/usr/bin/security", "delete-generic-password", "-a", "test", "-s", service).Run(); e != nil {
		t.Fatal("vault delete failed")
	}
	if e := exec.Command("/usr/bin/security", "find-generic-password", "-a", "test", "-s", service, "-w").Run(); e == nil {
		t.Fatal("credential remains after deletion")
	}
}
func TestEnvironmentPrecedence(t *testing.T) {
	t.Setenv("MIXRANK_API_KEY", "synthetic-env")
	v, source, e := Resolve()
	if e != nil || v != "synthetic-env" || source != "environment" {
		t.Fatal("environment override ignored")
	}
}
