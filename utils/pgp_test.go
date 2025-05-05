package utils

import (
	"testing"
)

func TestPGPEncryptDecrypt(t *testing.T) {
	data := "4532756279624982"

	encrypted, err := PGPEncrypt(data)
	if err != nil {
		t.Fatalf("PGPEncrypt failed: %v", err)
	}

	decrypted, err := PGPDecrypt(encrypted)
	if err != nil {
		t.Fatalf("PGPDecrypt failed: %v", err)
	}

	if decrypted != data {
		t.Errorf("Decrypted data does not match original: expected %s, got %s", data, decrypted)
	}
}
