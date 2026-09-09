// Package credentials owns only MixRank's OS vault entry. Host-managed plugin
// secrets arrive via environment variables; their credential files are untouched.
package credentials

import (
	"errors"
	"os"
	"strings"
)

const Service = "com.nkulavic.mixrank"
const Account = "default"

func Resolve() (string, string, error) {
	if v := os.Getenv("MIXRANK_API_KEY"); v != "" {
		return v, "environment", nil
	}
	v, e := read()
	if e != nil || v == "" {
		return "", "unavailable", errors.New("MixRank OS credential unavailable; run mixrank auth login or inject MIXRANK_API_KEY")
	}
	return v, "os-vault", nil
}
func Save(key string) error {
	if strings.TrimSpace(key) == "" || strings.ContainsAny(key, "\r\n\x00") {
		return errors.New("invalid API credential")
	}
	if e := save(key); e != nil {
		return errors.New("OS credential store rejected save; no plaintext fallback was written")
	}
	return nil
}
func Delete() error {
	if e := remove(); e != nil {
		return errors.New("OS credential store rejected delete or entry missing")
	}
	return nil
}
