package model

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"
)

// Envelope with signed user's data.
type Envelope struct {
	uuid    []byte
	expires time.Time
	payload []byte

	data      []byte
	signature []byte
}

// NewEnvelope constructor.
func NewEnvelope(
	secret string,
	useruuid string,
	payload string,
	expires time.Time,
) (*Envelope, error) {
	// validate input
	if len(secret) < 1 {
		return nil, errors.New("secret shouldn't be empty")
	}
	if len(payload) > 255 {
		return nil, errors.New("payload should be shorter then 256")
	}
	if _, err := uuid.Parse(useruuid); err != nil {
		return nil, fmt.Errorf("bad user uuid: %q", useruuid)
	}

	uuidLen := 36
	expiresLen := 12
	payloadLen := len([]byte(payload))

	b := make([]byte, uuidLen+expiresLen+1+payloadLen) // 1 byte for payloadLen

	// copy uuid bytes to b
	idx := copy(b, []byte(useruuid))

	// copy expires bytes to b
	seconds := []byte(fmt.Sprintf("%012d", expires.Unix()))
	idx += copy(b[idx:], seconds)

	// copy payload bytes to b
	b[idx] = byte(payloadLen)
	idx++
	idx += copy(b[idx:], []byte(payload))

	// sign with secret
	signature := sign(b, []byte(secret))

	signatureLen := len(signature)
	bSigned := make([]byte, len(b)+1+signatureLen)
	idx = copy(bSigned, b)
	bSigned[idx] = byte(signatureLen)
	idx++
	copy(bSigned[idx:], signature)

	envelope := &Envelope{
		uuid:    []byte(useruuid),
		expires: expires,
		payload: []byte(payload),

		data:      b,
		signature: signature,
	}

	return envelope, nil
}

// ParseEnvelope Envelope from string.
func ParseEnvelope(s string) (*Envelope, error) {
	b, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	// parse uuid
	idx := 0
	uuidLen := 36
	if len(b) < uuidLen {
		return nil, errors.New("bad input string")
	}
	useruuid := b[idx : idx+uuidLen]
	if _, err := uuid.Parse(string(useruuid)); err != nil {
		return nil, fmt.Errorf("bad uuid: %q", string(useruuid))
	}
	idx += uuidLen

	// parse expires
	expiresLen := 12
	if len(b) < (idx + expiresLen) {
		return nil, errors.New("bad input string")
	}
	seconds, err := strconv.ParseInt(string(b[idx:idx+expiresLen]), 10, 64)
	if err != nil {
		return nil, errors.New("bad input string")
	}
	expires := time.Unix(int64(seconds), 0)
	if expires.Before(time.Now().UTC()) {
		return nil, errors.New("expired")
	}
	idx += expiresLen

	// parse payload
	if len(b) < (idx + 1) {
		return nil, errors.New("bad input string")
	}
	payloadLen := int(b[idx])
	idx++
	if len(b) < (idx + payloadLen) {
		return nil, errors.New("bad input string")
	}
	payload := b[idx : idx+payloadLen]
	idx += payloadLen

	// parse signature
	if len(b) < (idx + 1) {
		return nil, errors.New("bad input string")
	}
	signatureLen := int(b[idx])
	idx++
	if len(b) < (idx + signatureLen) {
		return nil, errors.New("bad input string")
	}
	signature := b[idx : idx+signatureLen]
	idx += signatureLen
	b = make([]byte, uuidLen+expiresLen+1+payloadLen) // 1 byte for payloadLen
	idx = copy(b, useruuid)
	secondsBytes := []byte(fmt.Sprintf("%012d", expires.Unix()))
	idx += copy(b[idx:], secondsBytes)
	b[idx] = byte(payloadLen)
	idx++
	copy(b[idx:], payload)

	envelope := &Envelope{
		uuid:    useruuid,
		expires: expires,
		payload: payload,

		data:      b,
		signature: signature,
	}

	return envelope, nil
}

// String to apply Stringer interface.
func (envelope *Envelope) String() string {
	b := make([]byte, len(envelope.data)+1+len(envelope.signature))
	idx := copy(b, envelope.data)
	b[idx] = byte(len(envelope.signature))
	idx++
	copy(b[idx:], envelope.signature)
	return base64.URLEncoding.EncodeToString(b)
}

// UUID value.
func (envelope *Envelope) UUID() string {
	return string(envelope.uuid)
}

// Payload value.
func (envelope *Envelope) Payload() string {
	return string(envelope.payload)
}

// Expires time.
func (envelope *Envelope) Expires() time.Time {
	return envelope.expires
}

// CheckSignature validates Envelope signature.
func (envelope *Envelope) CheckSignature(secret string) bool {
	signature := sign(envelope.data, []byte(secret))
	return bytes.Compare(signature, envelope.signature) == 0
}

// sign signs data using key by calculating sha3-512 hash
// of their concatenation. Returns base64 url encoded string.
func sign(data, key []byte) []byte {
	b := make([]byte, len(data)+len(key))
	idx := copy(b, data)
	copy(b[idx:], key)
	h := sha3.Sum512(b)
	return h[:]
}
