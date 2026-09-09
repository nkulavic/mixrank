package mixrank

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// VerifyWebhook verifies the provider's raw API-key HMAC contract, including
// future timestamp rejection. Consumers must also deduplicate Webhook-Id values.
func VerifyWebhook(key string, headers http.Header, body []byte, now time.Time, tolerance time.Duration) error {
	if key == "" || tolerance <= 0 {
		return errors.New("webhook verification requires a key and positive tolerance")
	}
	id, ts := headers.Get("Webhook-Id"), headers.Get("Webhook-Timestamp")
	seconds, e := strconv.ParseInt(ts, 10, 64)
	if e != nil || id == "" {
		return errors.New("missing or invalid webhook headers")
	}
	delta := now.Sub(time.Unix(seconds, 0))
	if delta > tolerance || delta < -tolerance {
		return errors.New("webhook timestamp outside accepted window")
	}
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(id + "." + ts + "."))
	mac.Write(body)
	expected := mac.Sum(nil)
	for _, sig := range strings.Fields(headers.Get("Webhook-Signature")) {
		parts := strings.SplitN(sig, ",", 2)
		if len(parts) != 2 || parts[0] != "v1" {
			continue
		}
		value, e := base64.StdEncoding.DecodeString(parts[1])
		if e == nil && hmac.Equal(value, expected) {
			return nil
		}
	}
	return errors.New("invalid webhook signature")
}
