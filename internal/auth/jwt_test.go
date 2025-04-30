package auth_test

import (
	"testing"
	"time"

	"github.com/DryHop2/chirpy/internal/auth"
	"github.com/google/uuid"
)

func TestJWTCreateAndValidate(t *testing.T) {
	secret := "testsecret"
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

	validatedID, err := auth.ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("failed to validate JWT: %v", err)
	}

	if validatedID != userID {
		t.Fatalf("expected userID %s but got %s", userID, validatedID)
	}
}

func TestJWTExpired(t *testing.T) {
	secret := "testsecret"
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, secret, -time.Minute) // Already expired
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

	_, err = auth.ValidateJWT(token, secret)
	if err == nil {
		t.Fatal("expected error for expired JWT, got nil")
	}
}

func TestJWTInvalidSecret(t *testing.T) {
	secret := "testsecret"
	userID := uuid.New()

	token, err := auth.MakeJWT(userID, secret, time.Minute)
	if err != nil {
		t.Fatalf("failed to create JWT: %v", err)
	}

	_, err = auth.ValidateJWT(token, "wrongsecret")
	if err == nil {
		t.Fatal("expected error for invalid secret, got nil")
	}
}
