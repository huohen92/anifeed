package webui

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestPasswordHashingAndLegacyCompatibility(t *testing.T) {
	const password = "correct horse battery staple"
	salt := newPasswordSalt()
	hashed := hashPassword(password, salt)
	if !verifyPassword(password, salt, hashed) {
		t.Fatal("new password hash did not verify")
	}
	if verifyPassword("wrong", salt, hashed) {
		t.Fatal("wrong password verified")
	}

	legacy := sha256.Sum256([]byte(salt + ":" + password))
	if !verifyPassword(password, salt, hex.EncodeToString(legacy[:])) {
		t.Fatal("legacy sha256 password hash did not verify")
	}
}
