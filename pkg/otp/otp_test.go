package otp

import (
	"encoding/base32"
	"fmt"
	"testing"
	"time"
)

func TestComputeCode(t *testing.T) {
	var codeTests = []struct {
		secret string
		value  int64
		code   int
	}{
		{"TX3WVAPKQ3VOMR4WGOACTG24ZRUWGVYE", 1, 504196},
		{"TX3WVAPKQ3VOMR4WGOACTG24ZRUWGVYE", 5, 148870},
		{"TX3WVAPKQ3VOMR4WGOACTG24ZRUWGVYE", 10000, 215144},
	}

	for _, v := range codeTests {
		key, _ := base32.StdEncoding.DecodeString(v.secret)
		c, err := computeCode(key, v.value)
		if err != nil {
			t.Errorf("computeCode(%s, %06d) error: %v", v.secret, v.value, err)
		}

		if c != v.code {
			t.Errorf("computeCode(%s, %06d): got %06d expected %06d\n", v.secret, v.value, c, v.code)
		}
	}
}

func TestCheck(t *testing.T) {
	period := uint16(30)
	key, _ := base32.StdEncoding.DecodeString("TX3WVAPKQ3VOMR4WGOACTG24ZRUWGVYE")
	otp := New(key, period, 1)

	t0 := time.Now().UTC()
	tunix := t0.Unix() / int64(period)
	tunixBefore := tunix - 5
	tunixAfter := tunix + 35

	code, err := computeCode(otp.secret, tunix)
	if err != nil {
		t.Errorf("computeCode(%s, %d) error: %v", otp.secret, tunix, err)
	}
	codeBefore, err := computeCode(otp.secret, tunixBefore)
	if err != nil {
		t.Errorf("computeCode(%s, %d) error: %v", otp.secret, tunixBefore, err)
	}
	codeAfter, err := computeCode(otp.secret, tunixAfter)
	if err != nil {
		t.Errorf("computeCode(%s, %d) error: %v", otp.secret, tunixAfter, err)
	}

	var codeTests = []struct {
		code   string
		result bool
	}{
		{fmt.Sprintf("%06d", code), true},
		{fmt.Sprintf("%06d", codeBefore), false},
		{fmt.Sprintf("%06d", codeAfter), false},
	}

	for _, v := range codeTests {
		r, err := otp.Check(v.code)
		if err != nil {
			t.Errorf("error from code=%s: %v", v.code, err)
		}
		if r != v.result {
			t.Errorf("bad result from code=%s: got %t expected %t", v.code, r, v.result)
		}
	}
}

func TestProvisionURI(t *testing.T) {
	cases := []struct {
		user, issuer string
		out          string
	}{
		{"test", "", "otpauth://totp/test?secret=TX3WVAPKQ3VOMR4WGOACTG24ZRUWGVYE"},
		{"test", "Company", "otpauth://totp/Company:test?secret=TX3WVAPKQ3VOMR4WGOACTG24ZRUWGVYE&issuer=Company"},
	}

	for i, c := range cases {
		secret, _ := base32.StdEncoding.DecodeString("TX3WVAPKQ3VOMR4WGOACTG24ZRUWGVYE")
		otp := New(secret, 0, 1)
		got, _ := otp.ProvisionURI(c.user, c.issuer)
		if got != c.out {
			t.Errorf("%d: want %q, got %q", i, c.out, got)
		}
	}
}
