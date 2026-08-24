package cryptopriority2

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"regexp"
	"strings"
	"testing"
	"time"

	bip39 "github.com/luxfi/go-bip39"
	"golang.org/x/crypto/ssh"
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

func TestCertificateCSRAndSSHParsers(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	template := &x509.Certificate{SerialNumber: big.NewInt(7), Subject: pkix.Name{CommonName: "example.test"}, NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour), DNSNames: []string{"example.test"}, KeyUsage: x509.KeyUsageDigitalSignature}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	certificate := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
	info, _, err := ParseCertificate(certificate)
	if err != nil || info.Subject == "" || info.DNSNames[0] != "example.test" {
		t.Fatalf("certificate=%#v err=%v", info, err)
	}
	if _, err := PublicKeyFromCertificate(certificate); err != nil {
		t.Fatal(err)
	}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{Subject: pkix.Name{CommonName: "csr.example.test"}, DNSNames: []string{"csr.example.test"}}, key)
	if err != nil {
		t.Fatal(err)
	}
	csr, err := ParseCSR(string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDER})))
	if err != nil || !csr.SignatureValid {
		t.Fatalf("csr=%#v err=%v", csr, err)
	}
	sshKey, err := ssh.NewPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	host, err := ParseSSHHostKey("example.test " + string(ssh.MarshalAuthorizedKey(sshKey)))
	if err != nil || host.KeyType != "ssh-rsa" {
		t.Fatalf("host=%#v err=%v", host, err)
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
