//go:build !darwin

package credentials

import "github.com/zalando/go-keyring"

func save(service, account, key string) error      { return keyring.Set(service, account, key) }
func read(service, account string) (string, error) { return keyring.Get(service, account) }
func remove(service, account string) error         { return keyring.Delete(service, account) }
