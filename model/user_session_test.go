package model

import (
	"testing"
	"time"
)

func TestEnvelope(t *testing.T) {
	expires := time.Now().UTC().Add(time.Hour)
	secretKey := "app_secret_key"
	payload := []string{
		"VdjNo8P6dw2mVvqijt8b",
		"verify-email",
		"forgot-password",
	}

	for _, item := range payload {
		envelope, err := NewEnvelope(
			secretKey,
			"182717de-a19f-4d59-a40e-afd5fda729aa",
			item,
			expires,
		)

		parsed, err := ParseEnvelope(envelope.String())
		if err != nil {
			t.Error(err)
		}

		if !parsed.CheckSignature(secretKey) {
			t.Errorf("secret not valid")
		}
	}
}
