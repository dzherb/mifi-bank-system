package security

import (
	"github.com/dzherb/mifi-bank-system/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strconv"
	"testing"
	"time"
)

type token struct {
	iat float64
	sub string
	exp float64
}

func TestMain(m *testing.M) {
	Init(&config.Config{
		SecretKey:      "secret",
		AccessTokenTTL: time.Hour,
	})

	os.Exit(m.Run())
}

func TestTokenIssue(t *testing.T) {
	token, err := IssueAccessToken(23)
	if err != nil {
		t.Error(err)
		return
	}

	userID, err := ValidateToken(token)
	if err != nil {
		t.Error(err)
		return
	}

	if userID != 23 {
		t.Errorf("got sub %d, expected 23", userID)
	}
}

func TestTokenValidation(t *testing.T) {
	cases := []struct {
		token   token
		isValid bool
	}{
		{
			token: token{
				iat: timeToFloat64(time.Now()),
				sub: "1",
				exp: timeToFloat64(time.Now().Add(time.Second * 100)),
			},
			isValid: true,
		},
		{
			token: token{
				iat: timeToFloat64(time.Now()),
				sub: "1",
				exp: timeToFloat64(time.Now().Add(-time.Second * 100)),
			},
			isValid: false,
		},
		{
			token: token{
				iat: timeToFloat64(time.Now().Add(time.Second * 100)),
				sub: "1",
				exp: timeToFloat64(time.Now().Add(time.Second * 100)),
			},
			isValid: false,
		},
		{
			token: token{
				iat: timeToFloat64(time.Now()),
				sub: "not_a_number",
				exp: timeToFloat64(time.Now().Add(time.Second * 100)),
			},
			isValid: false,
		},
	}

	for _, c := range cases {
		tokenEncoded, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iat": c.token.iat,
			"sub": c.token.sub,
			"exp": c.token.exp,
		}).SignedString(secretKey)

		if err != nil {
			t.Error(err)
		}

		userID, err := ValidateToken(tokenEncoded)
		if err != nil {
			if c.isValid {
				t.Error("unexpected error:", err)
			}
			continue
		}

		if !c.isValid {
			t.Error("token is unexpectedly valid")
			continue
		}

		expectedID, err := strconv.Atoi(c.token.sub)
		if err != nil {
			t.Error(err)
			continue
		}

		if userID != expectedID {
			t.Errorf("got sub %d, want %d", userID, expectedID)
		}
	}
}
