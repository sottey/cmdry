package cryptopriority2

import (
	"regexp"
	"strings"
	"testing"
	"time"

	bip39 "github.com/luxfi/go-bip39"
)

func TestTokenAndULID(t *testing.T) {
	token, err := GenerateToken(16, "hex")
	if err != nil || len(token) != 32 {
		t.Fatalf("token=%q err=%v", token, err)
	}
	ulid, err := GenerateULID(time.UnixMilli(1724371200000))
	if err != nil || !regexp.MustCompile(`^[0-9ABCDEFGHJKMNPQRSTVWXYZ]{26}$`).MatchString(ulid) {
		t.Fatalf("ulid=%q err=%v", ulid, err)
	}
}

func TestEncryptionAndHMAC(t *testing.T) {
	payload, err := EncryptText("private text", "correct horse battery staple")
	if err != nil || !strings.HasPrefix(payload, encryptionVersion+":") {
		t.Fatalf("payload=%q err=%v", payload, err)
	}
	plaintext, err := DecryptText(payload, "correct horse battery staple")
	if err != nil || plaintext != "private text" {
		t.Fatalf("plain=%q err=%v", plaintext, err)
	}
	if _, err := DecryptText(payload, "wrong passphrase"); err == nil {
		t.Fatal("decrypted with wrong passphrase")
	}
	digest, err := HMAC("The quick brown fox jumps over the lazy dog", "key", "sha256")
	if err != nil || digest != "f7bc83f430538424b13298e6aa6fb143ef4d59a14946175997479dbc2d1a3cd8" {
		t.Fatalf("digest=%q err=%v", digest, err)
	}
}

func TestMnemonicAndPasswordStrength(t *testing.T) {
	mnemonic, err := GenerateMnemonic(12)
	if err != nil || len(strings.Fields(mnemonic)) != 12 || !bip39.IsMnemonicValid(mnemonic) {
		t.Fatalf("mnemonic=%q err=%v", mnemonic, err)
	}
	strength, err := AnalyzePassword("correct horse battery staple 2026!")
	if err != nil || strength.Score < 3 || strength.Entropy <= 0 {
		t.Fatalf("strength=%#v err=%v", strength, err)
	}
}
