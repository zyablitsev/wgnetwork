package otp

import (
	"fmt"
	"strings"

	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"net/url"
	"strconv"
	"time"
)

// OTP struct.
type OTP struct {
	secret []byte
	period uint16
	window int8
}

// New constructor.
func New(secret []byte, period uint16, window int8) *OTP {
	if period > 30 {
		period = 30
	}

	if window > 10 {
		window = 10
	} else if window < 1 {
		window = 1
	}

	return &OTP{secret: secret, period: period, window: window}
}

// ComputeCode for secret0.
func (c *OTP) ComputeCode() (string, error) {
	t := time.Now().UTC().Unix() / int64(c.period)

	code, err := computeCode(c.secret, t)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", code), nil
}

var errInvalidCode = errors.New("invalid code")

// Check a one-time-password against the given Config
// Returns true/false if the authentication was successful.
// Returns error if the password is incorrectly formatted (not a zero-padded 6 or non-zero-padded 8 digit number).
func (c *OTP) Check(password string) (bool, error) {
	if !(len(password) == 6 && password[0] >= '0' && password[0] <= '9') {
		return false, errInvalidCode
	}

	code, err := strconv.Atoi(password)
	if err != nil {
		return false, errInvalidCode
	}

	t := time.Now().UTC().Unix() / int64(c.period)
	for i := -c.window; i < c.window; i++ {
		tw := t - int64(i)
		computed, err := computeCode(c.secret, tw)
		if err != nil {
			return false, err
		}

		if computed == code {
			return true, nil
		}
	}

	return false, nil
}

// ProvisionURI generates a URI that can be turned into a QR code
// to configure a Google Authenticator mobile app.
func (c *OTP) ProvisionURI(user, issuer string) (string, string) {
	var buf strings.Builder

	buf.WriteString("otpauth://totp/")
	if issuer != "" {
		buf.WriteString(issuer)
		buf.WriteByte(':')
	}
	buf.WriteString(user)

	secret := base32.StdEncoding.EncodeToString(c.secret)
	buf.WriteByte('?')
	buf.WriteString("secret")
	buf.WriteByte('=')
	buf.WriteString(url.QueryEscape(secret))

	if issuer != "" {
		buf.WriteByte('&')
		buf.WriteString("issuer")
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(issuer))
	}

	return buf.String(), secret
}

func computeCode(secret []byte, value int64) (int, error) {
	hash := hmac.New(sha1.New, secret)
	err := binary.Write(hash, binary.BigEndian, value)
	if err != nil {
		return 0, err
	}
	h := hash.Sum(nil)

	offset := h[19] & 0x0f

	truncated := binary.BigEndian.Uint32(h[offset : offset+4])

	truncated &= 0x7fffffff
	code := truncated % 1000000

	return int(code), nil
}
