// Package credentials owns only the toolkit's provider entries in the OS vault.
// Host-managed plugin secrets arrive via environment variables; their
// credential files are untouched.
package credentials

import (
	"errors"
	"os"
	"strings"
)

const Service = "com.nkulavic.mixrank"
const Account = "default"
const PlacesService = "com.nkulavic.mixrank.google-places"
const PlacesAccount = "default"

func Resolve() (string, string, error) {
	if v := os.Getenv("MIXRANK_API_KEY"); v != "" {
		return v, "environment", nil
	}
	v, e := read(Service, Account)
	if e != nil || v == "" {
		return "", "unavailable", errors.New("MixRank OS credential unavailable; run mixrank auth login or inject MIXRANK_API_KEY")
	}
	return v, "os-vault", nil
}
func Save(key string) error {
	if strings.TrimSpace(key) == "" || strings.ContainsAny(key, "\r\n\x00") {
		return errors.New("invalid API credential")
	}
	if e := save(Service, Account, key); e != nil {
		return errors.New("OS credential store rejected save; no plaintext fallback was written")
	}
	return nil
}
func Delete() error {
	if e := remove(Service, Account); e != nil {
		return errors.New("OS credential store rejected delete or entry missing")
	}
	return nil
}

func ResolveGooglePlaces() (string, string, error) {
	for _, name := range []string{"GOOGLE_PLACES_API_KEY", "GOOGLE_MAPS_API_KEY"} {
		if v := os.Getenv(name); v != "" {
			return v, "environment", nil
		}
	}
	v, e := read(PlacesService, PlacesAccount)
	if e != nil || v == "" {
		return "", "unavailable", errors.New("Google Places OS credential unavailable; run mixrank auth google-places login or inject GOOGLE_PLACES_API_KEY")
	}
	return v, "os-vault", nil
}

func SaveGooglePlaces(key string) error {
	if strings.TrimSpace(key) == "" || strings.ContainsAny(key, "\r\n\x00") {
		return errors.New("invalid Google Places credential")
	}
	if e := save(PlacesService, PlacesAccount, key); e != nil {
		return errors.New("OS credential store rejected save; no plaintext fallback was written")
	}
	return nil
}

func DeleteGooglePlaces() error {
	if e := remove(PlacesService, PlacesAccount); e != nil {
		return errors.New("OS credential store rejected delete or entry missing")
	}
	return nil
}
