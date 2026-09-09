package mixrank

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ValidationPending records no email address, only the provider job identifier.
// An asynchronous one-off keeps the credential blocked until its signed callback
// is acknowledged; use bulk jobs when a callback receiver is not available.
type ValidationPending struct {
	ID         string    `json:"id"`
	AcceptedAt time.Time `json:"accepted_at"`
	Uncertain  bool      `json:"uncertain"`
}

var pendingMemory sync.Map

func (c *Client) coordinationID() string { return fmt.Sprintf("%x", sha256.Sum256([]byte(c.key))) }
func (c *Client) pendingPath() string {
	return filepath.Join(c.lockDir, c.coordinationID()+".pending.json")
}
func (c *Client) ValidationStatus() (*ValidationPending, error) {
	if c.volatile {
		if p, ok := pendingMemory.Load(c.coordinationID()); ok {
			v := p.(ValidationPending)
			return &v, nil
		}
		return nil, nil
	}
	b, e := os.ReadFile(c.pendingPath())
	if os.IsNotExist(e) {
		return nil, nil
	}
	if e != nil {
		return nil, errors.New("cannot read pending validation")
	}
	var p ValidationPending
	if json.Unmarshal(b, &p) != nil {
		return nil, errors.New("pending validation state unreadable; inspect before another one-off")
	}
	return &p, nil
}
func (c *Client) setPending(p ValidationPending) error {
	if c.volatile {
		pendingMemory.Store(c.coordinationID(), p)
		return nil
	}
	b, _ := json.Marshal(p)
	if e := os.WriteFile(c.pendingPath(), b, 0600); e != nil {
		return errors.New("cannot persist pending validation; stop one-off calls and coordinate with provider")
	}
	return nil
}
func (c *Client) clearPending() error {
	if c.volatile {
		pendingMemory.Delete(c.coordinationID())
		return nil
	}
	if e := os.Remove(c.pendingPath()); e != nil && !os.IsNotExist(e) {
		return errors.New("cannot clear pending validation")
	}
	return nil
}

// AcknowledgeValidationWebhook releases a one-off gate only for the matching,
// authenticated completion event. Callers must still deduplicate webhook IDs.
func (c *Client) AcknowledgeValidationWebhook(ctx context.Context, h http.Header, body []byte, now time.Time, tolerance time.Duration) error {
	if e := VerifyWebhook(c.key, h, body, now, tolerance); e != nil {
		return e
	}
	var event struct {
		ID   string          `json:"id"`
		Data json.RawMessage `json:"data"`
	}
	if json.Unmarshal(body, &event) != nil || event.ID == "" || event.ID != h.Get("Webhook-Id") || len(event.Data) == 0 {
		return errors.New("invalid validation completion event")
	}
	unlock, e := c.validationLock(ctx)
	if e != nil {
		return e
	}
	defer unlock()
	p, e := c.ValidationStatus()
	if e != nil {
		return e
	}
	if p == nil || p.ID != event.ID {
		return errors.New("webhook does not match pending one-off validation")
	}
	return c.clearPending()
}

// ClearValidation requires callers to have independently confirmed completion
// with their callback receiver/provider. It never submits a replacement request.
func (c *Client) ClearValidation(ctx context.Context, confirmed bool) error {
	if !confirmed {
		return errors.New("provider completion must be confirmed before clearing uncertain state")
	}
	unlock, e := c.validationLock(ctx)
	if e != nil {
		return e
	}
	defer unlock()
	return c.clearPending()
}
func (c *Client) validatePending() error {
	p, e := c.ValidationStatus()
	if e != nil {
		return e
	}
	if p != nil {
		return errors.New("prior one-off validation is pending or uncertain; acknowledge its signed webhook or confirm provider completion before clearing; use bulk jobs meanwhile")
	}
	return nil
}
