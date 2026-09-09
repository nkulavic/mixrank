//go:build !darwin

package credentials

import "github.com/zalando/go-keyring"

func save(key string) error { return keyring.Set(Service, Account, key) }
func read() (string, error) { return keyring.Get(Service, Account) }
func remove() error         { return keyring.Delete(Service, Account) }
