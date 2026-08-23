package cryptopriority2

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
)

func RunTokenGenerator() {
	cmdry.Run(plugin("token-generator", "Token Generator", "Generate bounded cryptographically secure tokens locally.", []string{"token", "random", "secure", "secret", "generator"}, tokenForm, tokenRun))
}
func RunULIDGenerator() {
	cmdry.Run(plugin("ulid-generator", "ULID Generator", "Generate sortable ULID values locally.", []string{"ulid", "identifier", "sortable", "random", "generator"}, ulidForm, ulidRun))
}
func RunEncryptDecryptText() {
	cmdry.Run(plugin("encrypt-decrypt-text", "Encrypt and Decrypt Text", "Encrypt or decrypt pasted text with a local passphrase; no input is stored.", []string{"encrypt", "decrypt", "aes", "passphrase", "scrypt"}, encryptionForm, encryptionRun))
}
func RunBIP39PassphraseGenerator() {
	cmdry.Run(plugin("bip39-passphrase-generator", "BIP39 Passphrase Generator", "Generate a BIP39 mnemonic locally; treat it as a secret.", []string{"bip39", "mnemonic", "seed phrase", "wallet", "passphrase"}, mnemonicForm, mnemonicRun))
}
func RunHMACGenerator() {
	cmdry.Run(plugin("hmac-generator", "HMAC Generator", "Create an HMAC digest from pasted text and a secret key without retaining either value.", []string{"hmac", "sha256", "sha512", "digest", "secret"}, hmacForm, hmacRun))
}
func RunRSAKeyPairGenerator() {
	cmdry.Run(plugin("rsa-key-pair-generator", "RSA Key Pair Generator", "Generate a bounded RSA key pair locally and provide browser downloads only.", []string{"rsa", "public key", "private key", "pem", "key pair"}, rsaForm, rsaRun))
}
func RunPasswordStrengthAnalyzer() {
	cmdry.Run(plugin("password-strength-analyzer", "Password Strength Analyzer", "Estimate password strength locally without storing or returning the password.", []string{"password", "strength", "entropy", "security", "analyze"}, strengthForm, strengthRun))
}

func plugin(id, name, description string, terms []string, overview, action cmdry.Handler) cmdry.Plugin {
	return cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey." + id, Name: name, Version: "0.1.0", Description: description, Category: "crypto", Icon: "lock", SearchTerms: terms, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "run", Name: "Run", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "run": action}}
}

func tokenForm(cmdry.Request) (cmdry.View, error) {
	return form("Token Generator", "Tokens are generated with cryptographic randomness and are shown only in this response.", "Generate token", []cmdry.Field{{Name: "bytes", Label: "Random bytes", Type: "number", Value: "32", Min: "16", Max: "64", Required: true}, {Name: "encoding", Label: "Encoding", Type: "select", Value: "base64url", Options: []cmdry.Option{{Value: "base64url", Label: "Base64 URL"}, {Value: "hex", Label: "Hexadecimal"}}}}), nil
}
func tokenRun(request cmdry.Request) (cmdry.View, error) {
	bytes, err := integer(request, "bytes")
	if err != nil {
		return cmdry.View{}, err
	}
	token, err := GenerateToken(bytes, fmt.Sprint(request.Params["encoding"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return codeResult("Generated token", "Copy this token now. Cmdry does not retain it.", token, "Generate another token"), nil
}

func ulidForm(cmdry.Request) (cmdry.View, error) {
	return oneClick("ULID Generator", "Generate a time-sortable ULID using the current time and cryptographic randomness.", "Generate ULID"), nil
}
func ulidRun(cmdry.Request) (cmdry.View, error) {
	value, err := GenerateULID(time.Now())
	if err != nil {
		return cmdry.View{}, err
	}
	return codeResult("Generated ULID", "The timestamp is encoded in the first 10 characters; the remaining characters are random.", value, "Generate another ULID"), nil
}

func encryptionForm(cmdry.Request) (cmdry.View, error) {
	return form("Encrypt and Decrypt Text", "Uses AES-256-GCM with a key derived locally from your passphrase. Forgetting the passphrase makes the data unrecoverable.", "Run encryption", []cmdry.Field{{Name: "operation", Label: "Operation", Type: "select", Value: "encrypt", Options: []cmdry.Option{{Value: "encrypt", Label: "Encrypt text"}, {Value: "decrypt", Label: "Decrypt Cmdry payload"}}}, {Name: "input", Label: "Text or encrypted payload", Type: "textarea", Required: true}, {Name: "passphrase", Label: "Passphrase", Type: "password", Required: true}}), nil
}
func encryptionRun(request cmdry.Request) (cmdry.View, error) {
	input, passphrase := fmt.Sprint(request.Params["input"]), fmt.Sprint(request.Params["passphrase"])
	if fmt.Sprint(request.Params["operation"]) == "decrypt" {
		output, err := DecryptText(input, passphrase)
		if err != nil {
			return cmdry.View{}, err
		}
		return codeResult("Decrypted text", "The text was decrypted locally and is not retained.", output, "Decrypt another value"), nil
	}
	output, err := EncryptText(input, passphrase)
	if err != nil {
		return cmdry.View{}, err
	}
	return codeResult("Encrypted payload", "Copy this payload and save the passphrase separately. Cmdry does not retain either value.", output, "Encrypt another value"), nil
}

func mnemonicForm(cmdry.Request) (cmdry.View, error) {
	return form("BIP39 Passphrase Generator", "A mnemonic can control funds or identities. Store it securely; Cmdry does not retain it.", "Generate mnemonic", []cmdry.Field{{Name: "words", Label: "Words", Type: "select", Value: "12", Options: []cmdry.Option{{Value: "12", Label: "12 words"}, {Value: "15", Label: "15 words"}, {Value: "18", Label: "18 words"}, {Value: "21", Label: "21 words"}, {Value: "24", Label: "24 words"}}}}), nil
}
func mnemonicRun(request cmdry.Request) (cmdry.View, error) {
	words, err := integer(request, "words")
	if err != nil {
		return cmdry.View{}, err
	}
	mnemonic, err := GenerateMnemonic(words)
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Generated BIP39 mnemonic", Components: []cmdry.Component{{Type: "alert", Level: "warning", Title: "Treat this as a secret", Message: "Anyone with this mnemonic may control the account or wallet derived from it. Cmdry does not retain it."}, {Type: "code", Title: "BIP39 mnemonic", Text: mnemonic}, restart("Generate another mnemonic")}}, nil
}

func hmacForm(cmdry.Request) (cmdry.View, error) {
	return form("HMAC Generator", "Create a keyed digest locally. The message and secret are not retained.", "Generate HMAC", []cmdry.Field{{Name: "algorithm", Label: "Algorithm", Type: "select", Value: "sha256", Options: []cmdry.Option{{Value: "sha256", Label: "HMAC-SHA-256"}, {Value: "sha512", Label: "HMAC-SHA-512"}}}, {Name: "message", Label: "Message", Type: "textarea", Required: true}, {Name: "secret", Label: "Secret key", Type: "password", Required: true}}), nil
}
func hmacRun(request cmdry.Request) (cmdry.View, error) {
	digest, err := HMAC(fmt.Sprint(request.Params["message"]), fmt.Sprint(request.Params["secret"]), fmt.Sprint(request.Params["algorithm"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return codeResult("HMAC digest", "The digest was created locally; the message and secret are not retained.", digest, "Generate another HMAC"), nil
}

func rsaForm(cmdry.Request) (cmdry.View, error) {
	return form("RSA Key Pair Generator", "Generate a private/public RSA key pair locally. The private key is provided only as a download; store it securely.", "Generate RSA key pair", []cmdry.Field{{Name: "bits", Label: "Key size", Type: "select", Value: "2048", Options: []cmdry.Option{{Value: "2048", Label: "2048 bits"}, {Value: "3072", Label: "3072 bits"}}}}), nil
}
func rsaRun(request cmdry.Request) (cmdry.View, error) {
	bits, err := integer(request, "bits")
	if err != nil {
		return cmdry.View{}, err
	}
	privateKey, publicKey, err := GenerateRSAKeyPair(bits)
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "RSA key pair", Components: []cmdry.Component{{Type: "alert", Level: "warning", Title: "Private key download", Message: "Save the private key securely. It is not retained by Cmdry."}, {Type: "download", Filename: "cmdry-rsa-private-key.pem", MIMEType: "application/x-pem-file", Content: base64.StdEncoding.EncodeToString([]byte(privateKey))}, {Type: "download", Filename: "cmdry-rsa-public-key.pem", MIMEType: "application/x-pem-file", Content: base64.StdEncoding.EncodeToString([]byte(publicKey))}, restart("Generate another key pair")}}, nil
}

func strengthForm(cmdry.Request) (cmdry.View, error) {
	return form("Password Strength Analyzer", "Estimates strength locally without storing or returning the password. This is guidance, not a guarantee.", "Analyze password", []cmdry.Field{{Name: "password", Label: "Password", Type: "password", Required: true}}), nil
}
func strengthRun(request cmdry.Request) (cmdry.View, error) {
	strength, err := AnalyzePassword(fmt.Sprint(request.Params["password"]))
	if err != nil {
		return cmdry.View{}, err
	}
	rows := []map[string]any{{"field": "Estimate", "value": strength.Label}, {"field": "Score", "value": fmt.Sprintf("%d / 4", strength.Score)}, {"field": "Approximate entropy", "value": fmt.Sprintf("%.0f bits", strength.Entropy)}, {"field": "Guidance", "value": strings.Join(strength.Feedback, " ")}}
	return tableResult("Password strength", "The password itself is not shown or retained.", "password-strength", rows, "Analyze another password"), nil
}

func form(title, description, submit string, fields []cmdry.Field) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "form", Title: title, Description: description + " Your data never leaves this device.", Action: "run", Submit: submit, Fields: fields}}}
}
func oneClick(title, description, button string) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "text", Title: title, Text: description + " Your data never leaves this device."}, {Type: "actions", Actions: []cmdry.Action{{ID: "run", Name: button, Method: "write"}}}}}
}
func codeResult(title, message, value, again string) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "Complete", Message: message}, {Type: "code", Title: title, Text: value}, restart(again)}}
}
func tableResult(title, message, id string, rows []map[string]any, again string) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "alert", Level: "info", Title: "Local estimate", Message: message}, {Type: "table", ID: id, Columns: []cmdry.Column{{Key: "field", Label: "Field"}, {Key: "value", Label: "Value"}}, Rows: rows}, restart(again)}}
}
func restart(name string) cmdry.Component {
	return cmdry.Component{Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: name}}}
}
func integer(request cmdry.Request, key string) (int, error) {
	value, err := strconv.Atoi(strings.TrimSpace(fmt.Sprint(request.Params[key])))
	if err != nil {
		return 0, fmt.Errorf("%s must be a whole number", strings.ReplaceAll(key, "_", " "))
	}
	return value, nil
}
