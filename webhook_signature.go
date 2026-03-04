package postproxy

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// VerifyWebhookSignature verifies the signature of a webhook payload.
func VerifyWebhookSignature(payload, signatureHeader, secret string) bool {
	parts := make(map[string]string)
	for _, p := range strings.Split(signatureHeader, ",") {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) == 2 {
			parts[kv[0]] = kv[1]
		}
	}

	timestamp := parts["t"]
	expected := parts["v1"]
	if timestamp == "" || expected == "" {
		return false
	}

	signedPayload := timestamp + "." + payload
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signedPayload))
	computed := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(computed), []byte(expected))
}
