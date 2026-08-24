// Package cryptopriority2 provides bounded, local-only cryptographic utilities.
package cryptopriority2

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"hash"
	"math"
	"math/big"
	"strings"
	"time"
	"unicode"

	bip39 "github.com/luxfi/go-bip39"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/ssh"
)

const encryptionVersion = "cmdry:v1"

type CertificateInfo struct {
	Subject, Issuer, Serial, NotBefore, NotAfter, PublicKeyAlgorithm, SignatureAlgorithm, SHA256 string
	DNSNames                                                                                     []string
	ExtensionCount                                                                               int
}

func decodePEMOrDER(input string, allowed ...string) ([]byte, error) {
	input = strings.TrimSpace(input)
	if len(input) > 1<<20 {
		return nil, fmt.Errorf("input must not exceed 1 MiB")
	}
	if block, rest := pem.Decode([]byte(input)); block != nil {
		if strings.TrimSpace(string(rest)) != "" {
			return nil, fmt.Errorf("enter exactly one PEM block")
		}
		for _, kind := range allowed {
			if block.Type == kind {
				return block.Bytes, nil
			}
		}
		return nil, fmt.Errorf("unsupported PEM block type %q", block.Type)
	}
	der, err := base64.StdEncoding.DecodeString(strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, input))
	if err != nil {
		return nil, fmt.Errorf("enter PEM or Base64-encoded DER")
	}
	return der, nil
}
func ParseCertificate(input string) (CertificateInfo, *x509.Certificate, error) {
	der, err := decodePEMOrDER(input, "CERTIFICATE", "X509 CERTIFICATE")
	if err != nil {
		return CertificateInfo{}, nil, err
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return CertificateInfo{}, nil, fmt.Errorf("parse certificate: %w", err)
	}
	digest := sha256.Sum256(der)
	info := CertificateInfo{Subject: cert.Subject.String(), Issuer: cert.Issuer.String(), Serial: cert.SerialNumber.String(), NotBefore: cert.NotBefore.UTC().Format(time.RFC3339), NotAfter: cert.NotAfter.UTC().Format(time.RFC3339), PublicKeyAlgorithm: cert.PublicKeyAlgorithm.String(), SignatureAlgorithm: cert.SignatureAlgorithm.String(), SHA256: strings.ToUpper(hex.EncodeToString(digest[:])), DNSNames: cert.DNSNames, ExtensionCount: len(cert.Extensions)}
	return info, cert, nil
}
func PublicKeyFromCertificate(input string) (string, error) {
	_, cert, err := ParseCertificate(input)
	if err != nil {
		return "", err
	}
	der, err := x509.MarshalPKIXPublicKey(cert.PublicKey)
	if err != nil {
		return "", fmt.Errorf("encode public key: %w", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})), nil
}

type CSRInfo struct {
	Subject, SignatureAlgorithm, PublicKeyAlgorithm string
	DNSNames                                        []string
	SignatureValid                                  bool
}

func ParseCSR(input string) (CSRInfo, error) {
	der, err := decodePEMOrDER(input, "CERTIFICATE REQUEST", "NEW CERTIFICATE REQUEST")
	if err != nil {
		return CSRInfo{}, err
	}
	request, err := x509.ParseCertificateRequest(der)
	if err != nil {
		return CSRInfo{}, fmt.Errorf("parse CSR: %w", err)
	}
	err = request.CheckSignature()
	return CSRInfo{Subject: request.Subject.String(), SignatureAlgorithm: request.SignatureAlgorithm.String(), PublicKeyAlgorithm: request.PublicKeyAlgorithm.String(), DNSNames: request.DNSNames, SignatureValid: err == nil}, err
}

type SSHHostKeyInfo struct{ Marker, Hosts, KeyType, Fingerprint, Comment string }

func ParseSSHHostKey(input string) (SSHHostKeyInfo, error) {
	fields := strings.Fields(strings.TrimSpace(input))
	if len(fields) < 3 || strings.Contains(strings.ToUpper(input), "PRIVATE KEY") {
		return SSHHostKeyInfo{}, fmt.Errorf("enter a public known_hosts or authorized-key line, not a private key")
	}
	info := SSHHostKeyInfo{}
	index := 0
	if strings.HasPrefix(fields[index], "@") {
		info.Marker = fields[index]
		index++
	}
	if len(fields) < index+3 {
		return SSHHostKeyInfo{}, fmt.Errorf("enter a complete public host-key line")
	}
	info.Hosts = fields[index]
	index++
	key, comment, _, _, err := ssh.ParseAuthorizedKey([]byte(strings.Join(fields[index:], " ")))
	if err != nil {
		return SSHHostKeyInfo{}, fmt.Errorf("parse SSH public key: %w", err)
	}
	info.KeyType = key.Type()
	info.Fingerprint = ssh.FingerprintSHA256(key)
	info.Comment = comment
	return info, nil
}

func GenerateToken(bytes int, encoding string) (string, error) {
	if bytes < 16 || bytes > 64 {
		return "", fmt.Errorf("token size must be between 16 and 64 bytes")
	}
	value := make([]byte, bytes)
	if _, err := rand.Read(value); err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	switch encoding {
	case "base64url":
		return base64.RawURLEncoding.EncodeToString(value), nil
	case "hex":
		return hex.EncodeToString(value), nil
	default:
		return "", fmt.Errorf("unsupported token encoding")
	}
}

// GenerateULID creates a sortable ULID using the current timestamp and 80 random bits.
func GenerateULID(now time.Time) (string, error) {
	if now.UnixMilli() < 0 || now.UnixMilli() >= 1<<48 {
		return "", fmt.Errorf("current time is outside the ULID range")
	}
	value := make([]byte, 16)
	milliseconds := uint64(now.UnixMilli())
	for index := 5; index >= 0; index-- {
		value[index] = byte(milliseconds)
		milliseconds >>= 8
	}
	if _, err := rand.Read(value[6:]); err != nil {
		return "", fmt.Errorf("generate ULID randomness: %w", err)
	}
	return crockfordBase32(value), nil
}

func crockfordBase32(value []byte) string {
	const alphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	number := new(big.Int).SetBytes(value)
	base, remainder := big.NewInt(32), new(big.Int)
	encoded := make([]byte, 26)
	for index := len(encoded) - 1; index >= 0; index-- {
		number.DivMod(number, base, remainder)
		encoded[index] = alphabet[remainder.Int64()]
	}
	return string(encoded)
}

// EncryptText uses scrypt-derived AES-256-GCM and returns a versioned portable payload.
func EncryptText(plaintext, passphrase string) (string, error) {
	if plaintext == "" || len(plaintext) > 1<<20 {
		return "", fmt.Errorf("text must contain between 1 and 1,048,576 bytes")
	}
	if len(passphrase) < 8 || len(passphrase) > 1024 {
		return "", fmt.Errorf("passphrase must contain between 8 and 1,024 characters")
	}
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate encryption salt: %w", err)
	}
	key, err := deriveKey(passphrase, salt)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create AES cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create AES-GCM: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("generate encryption nonce: %w", err)
	}
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	return strings.Join([]string{encryptionVersion, base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(nonce), base64.RawStdEncoding.EncodeToString(ciphertext)}, ":"), nil
}

func DecryptText(payload, passphrase string) (string, error) {
	if len(payload) == 0 || len(payload) > 2<<20 {
		return "", fmt.Errorf("encrypted payload must contain between 1 and 2,097,152 bytes")
	}
	if len(passphrase) < 8 || len(passphrase) > 1024 {
		return "", fmt.Errorf("passphrase must contain between 8 and 1,024 characters")
	}
	parts := strings.Split(payload, ":")
	if len(parts) != 5 || parts[0]+":"+parts[1] != encryptionVersion {
		return "", fmt.Errorf("encrypted payload is not a supported Cmdry v1 value")
	}
	salt, saltErr := base64.RawStdEncoding.DecodeString(parts[2])
	nonce, nonceErr := base64.RawStdEncoding.DecodeString(parts[3])
	ciphertext, cipherErr := base64.RawStdEncoding.DecodeString(parts[4])
	if saltErr != nil || nonceErr != nil || cipherErr != nil || len(salt) != 16 {
		return "", fmt.Errorf("encrypted payload is malformed")
	}
	key, err := deriveKey(passphrase, salt)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create AES cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil || len(nonce) != gcm.NonceSize() {
		return "", fmt.Errorf("encrypted payload has an invalid nonce")
	}
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decrypt payload: incorrect passphrase or modified data")
	}
	return string(plaintext), nil
}

func deriveKey(passphrase string, salt []byte) ([]byte, error) {
	key, err := scrypt.Key([]byte(passphrase), salt, 32768, 8, 1, 32)
	if err != nil {
		return nil, fmt.Errorf("derive encryption key: %w", err)
	}
	return key, nil
}

func GenerateMnemonic(words int) (string, error) {
	bitsByWords := map[int]int{12: 128, 15: 160, 18: 192, 21: 224, 24: 256}
	bits, ok := bitsByWords[words]
	if !ok {
		return "", fmt.Errorf("choose 12, 15, 18, 21, or 24 words")
	}
	entropy, err := bip39.NewEntropy(bits)
	if err != nil {
		return "", fmt.Errorf("generate BIP39 entropy: %w", err)
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("generate BIP39 mnemonic: %w", err)
	}
	return mnemonic, nil
}

func HMAC(message, secret, algorithm string) (string, error) {
	if message == "" || len(message) > 1<<20 {
		return "", fmt.Errorf("message must contain between 1 and 1,048,576 bytes")
	}
	if secret == "" || len(secret) > 4096 {
		return "", fmt.Errorf("secret must contain between 1 and 4,096 bytes")
	}
	var constructor func() hash.Hash
	switch algorithm {
	case "sha256":
		constructor = sha256.New
	case "sha512":
		constructor = sha512.New
	default:
		return "", fmt.Errorf("unsupported HMAC algorithm")
	}
	mac := hmac.New(constructor, []byte(secret))
	_, _ = mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil)), nil
}

func GenerateRSAKeyPair(bits int) (privatePEM, publicPEM string, err error) {
	if bits != 2048 && bits != 3072 {
		return "", "", fmt.Errorf("RSA key size must be 2048 or 3072 bits")
	}
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", fmt.Errorf("generate RSA key pair: %w", err)
	}
	privatePEM = string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}))
	publicDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("encode RSA public key: %w", err)
	}
	publicPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicDER}))
	return privatePEM, publicPEM, nil
}

type PasswordStrength struct {
	Score    int
	Label    string
	Entropy  float64
	Feedback []string
}

// AnalyzePassword estimates strength without storing or returning the password.
func AnalyzePassword(password string) (PasswordStrength, error) {
	length := len([]rune(password))
	if length == 0 || length > 1024 {
		return PasswordStrength{}, fmt.Errorf("password must contain between 1 and 1,024 characters")
	}
	classes, pool := 0, 0.0
	hasLower, hasUpper, hasDigit, hasSymbol := false, false, false, false
	for _, character := range password {
		switch {
		case unicode.IsLower(character):
			hasLower = true
		case unicode.IsUpper(character):
			hasUpper = true
		case unicode.IsDigit(character):
			hasDigit = true
		default:
			hasSymbol = true
		}
	}
	if hasLower {
		classes++
		pool += 26
	}
	if hasUpper {
		classes++
		pool += 26
	}
	if hasDigit {
		classes++
		pool += 10
	}
	if hasSymbol {
		classes++
		pool += 33
	}
	entropy := float64(length) * math.Log2(pool)
	score := 0
	if length >= 8 {
		score++
	}
	if length >= 12 {
		score++
	}
	if classes >= 3 {
		score++
	}
	if length >= 16 && classes >= 4 {
		score++
	}
	feedback := []string{}
	lower := strings.ToLower(password)
	if containsCommonPattern(lower) {
		score = max(0, score-2)
		entropy *= 0.45
		feedback = append(feedback, "Avoid common words, keyboard sequences, and repeated characters.")
	}
	if length < 12 {
		feedback = append(feedback, "Use at least 12 characters.")
	}
	if classes < 3 {
		feedback = append(feedback, "Mix uppercase, lowercase, numbers, and symbols.")
	}
	labels := []string{"Very weak", "Weak", "Fair", "Strong", "Very strong"}
	if len(feedback) == 0 {
		feedback = append(feedback, "This estimate is stronger when the password is unique and not reused.")
	}
	return PasswordStrength{Score: score, Label: labels[score], Entropy: entropy, Feedback: feedback}, nil
}

func containsCommonPattern(value string) bool {
	for _, pattern := range []string{"password", "qwerty", "123456", "letmein", "admin", "welcome", "abcdef"} {
		if strings.Contains(value, pattern) {
			return true
		}
	}
	for index := 2; index < len(value); index++ {
		if value[index] == value[index-1] && value[index] == value[index-2] {
			return true
		}
	}
	return false
}
