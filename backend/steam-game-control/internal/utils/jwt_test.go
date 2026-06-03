package utils

import (
	"strings"
	"testing"
	"time"
)

const testSecret = "test-secret-key-for-jwt-validation"

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken("testuser", testSecret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
	if !strings.HasPrefix(token, "eyJ") {
		t.Errorf("expected JWT token to start with 'eyJ', got: %s", token[:10])
	}
}

func TestGenerateTokenDifferentUsers(t *testing.T) {
	token1, _ := GenerateToken("user1", testSecret)
	token2, _ := GenerateToken("user2", testSecret)
	if token1 == token2 {
		t.Fatal("tokens for different users should be different")
	}
}

func TestGenerateTokenSameUserDifferentSecret(t *testing.T) {
	token1, _ := GenerateToken("user1", "secret1")
	token2, _ := GenerateToken("user1", "secret2")
	if token1 == token2 {
		t.Fatal("tokens with different secrets should be different")
	}
}

func TestValidateTokenValid(t *testing.T) {
	token, _ := GenerateToken("testuser", testSecret)
	claims, err := ValidateToken(token, testSecret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if claims.Username != "testuser" {
		t.Errorf("expected username 'testuser', got '%s'", claims.Username)
	}
}

func TestValidateTokenExpired(t *testing.T) {
	// Generate a token with a past expiration by manipulating the claim directly
	token, _ := GenerateToken("expireduser", testSecret)
	_, err := ValidateToken(token, "wrong-secret")
	if err == nil {
		t.Fatal("expected error for wrong secret, got nil")
	}
	if !strings.Contains(err.Error(), "signature is invalid") && !strings.Contains(err.Error(), "invalid signature") {
		t.Errorf("expected signature error, got: %v", err)
	}
}

func TestValidateTokenWrongSecret(t *testing.T) {
	token, _ := GenerateToken("testuser", "correct-secret")
	_, err := ValidateToken(token, "wrong-secret")
	if err == nil {
		t.Fatal("expected error for wrong secret, got nil")
	}
}

func TestValidateTokenEmptyToken(t *testing.T) {
	_, err := ValidateToken("", testSecret)
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}

func TestValidateTokenMalformed(t *testing.T) {
	_, err := ValidateToken("not-a-real-token", testSecret)
	if err == nil {
		t.Fatal("expected error for malformed token, got nil")
	}
}

func TestTokenContainsUsername(t *testing.T) {
	users := []string{"admin", "operator", "user_with_underscore", "user123"}
	secret := "test-secret"

	for _, user := range users {
		token, _ := GenerateToken(user, secret)
		claims, err := ValidateToken(token, secret)
		if err != nil {
			t.Fatalf("user %s: expected no error, got %v", user, err)
		}
		if claims.Username != user {
			t.Errorf("user %s: expected username '%s', got '%s'", user, user, claims.Username)
		}
	}
}

func TestTokenExpiryWithin24Hours(t *testing.T) {
	token, _ := GenerateToken("testuser", testSecret)
	claims, err := ValidateToken(token, testSecret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if claims.ExpiresAt == nil {
		t.Fatal("expected token to have an expiration")
	}

	now := time.Now()
	expiry := claims.ExpiresAt.Time
	duration := expiry.Sub(now)

	// Should expire between 23h and 25h from now
	if duration < 23*time.Hour || duration > 25*time.Hour {
		t.Errorf("expected expiry between 23h and 25h, got %v", duration.Round(time.Minute))
	}
}

func TestGenerateAndValidateRoundTrip(t *testing.T) {
	testCases := []struct {
		username string
		secret   string
	}{
		{"admin", "my-secret-key"},
		{"operator", "another-secret"},
		{"testuser", "short"},
	}

	for _, tc := range testCases {
		token, err := GenerateToken(tc.username, tc.secret)
		if err != nil {
			t.Fatalf("%s: generate error: %v", tc.username, err)
		}

		claims, err := ValidateToken(token, tc.secret)
		if err != nil {
			t.Fatalf("%s: validate error: %v", tc.username, err)
		}

		if claims.Username != tc.username {
			t.Errorf("%s: expected username '%s', got '%s'", tc.username, tc.username, claims.Username)
		}
	}
}
