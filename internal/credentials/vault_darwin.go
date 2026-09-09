//go:build darwin

package credentials

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

// security's interactive command interpreter keeps the password off argv.
// Its output is captured and never surfaced because diagnostics can echo input.
func save(key string) error {
	if strings.ContainsAny(key, "\"'\\") {
		return errors.New("credential contains characters unsupported by the Keychain command interpreter; use runtime environment injection")
	}
	cmd := exec.Command("/usr/bin/security", "-i")
	cmd.Stdin = strings.NewReader("add-generic-password -U -a " + Account + " -s " + Service + " -w \"" + key + "\"\n")
	out, e := cmd.CombinedOutput()
	if e != nil || bytes.Contains(out, []byte("SecKeychain")) || bytes.Contains(out, []byte("Error:")) {
		return errors.New("keychain save failed")
	}
	v, e := read()
	if e != nil || v != key {
		return errors.New("keychain save not verified")
	}
	return nil
}
func read() (string, error) {
	b, e := exec.Command("/usr/bin/security", "find-generic-password", "-a", Account, "-s", Service, "-w").Output()
	return strings.TrimSuffix(string(b), "\n"), e
}
func remove() error {
	return exec.Command("/usr/bin/security", "delete-generic-password", "-a", Account, "-s", Service).Run()
}
