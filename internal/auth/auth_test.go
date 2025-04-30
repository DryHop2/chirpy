package auth

import (
	"testing"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "supersecure123"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	if hashed == password {
		t.Errorf("hashed password should not be equal to raw password")
	}

	// Test correct password
	if err := CheckPasswordHash(hashed, password); err != nil {
		t.Errorf("expected password to match, got error: %v", err)
	}

	// Test incorrect password
	if err := CheckPasswordHash(hashed, "wrongpassword"); err == nil {
		t.Error("expected error for wrong password, got nil")
	}
}
