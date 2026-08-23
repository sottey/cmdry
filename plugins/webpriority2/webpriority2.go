package webpriority2

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func BasicAuth(user, password string) (string, error) {
	if user == "" || strings.ContainsAny(user, "\r\n:") {
		return "", fmt.Errorf("username is required and cannot contain colon or newlines")
	}
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+password)), nil
}
func Entity(input, op string) (string, error) {
	if len(input) > 1<<20 {
		return "", fmt.Errorf("input exceeds 1 MiB")
	}
	if op == "encode" {
		return html.EscapeString(input), nil
	}
	if op == "decode" {
		return html.UnescapeString(input), nil
	}
	return "", fmt.Errorf("choose encode or decode")
}

type JWT struct {
	Header, Claims string
	Signature      string
}

func InspectJWT(input string) (JWT, error) {
	parts := strings.Split(strings.TrimSpace(input), ".")
	if len(parts) != 3 {
		return JWT{}, fmt.Errorf("JWT must contain three dot-separated parts")
	}
	decode := func(value string) (string, error) {
		raw, err := base64.RawURLEncoding.DecodeString(value)
		if err != nil {
			return "", err
		}
		var v any
		if json.Unmarshal(raw, &v) != nil {
			return "", fmt.Errorf("JWT segment is not JSON")
		}
		pretty, _ := json.MarshalIndent(v, "", "  ")
		return string(pretty), nil
	}
	header, err := decode(parts[0])
	if err != nil {
		return JWT{}, fmt.Errorf("header: %w", err)
	}
	claims, err := decode(parts[1])
	if err != nil {
		return JWT{}, fmt.Errorf("claims: %w", err)
	}
	return JWT{header, claims, parts[2]}, nil
}
func OTP(secret string, now time.Time) (string, error) {
	secret = strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(secret), " ", ""))
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil || len(key) == 0 {
		return "", fmt.Errorf("enter a valid Base32 TOTP secret")
	}
	counter := uint64(now.Unix() / 30)
	message := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		message[i] = byte(counter)
		counter >>= 8
	}
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	sum := mac.Sum(nil)
	offset := sum[len(sum)-1] & 15
	code := (int(sum[offset])&127)<<24 | (int(sum[offset+1])&255)<<16 | (int(sum[offset+2])&255)<<8 | (int(sum[offset+3]) & 255)
	return fmt.Sprintf("%06d", code%1000000), nil
}
func SafeLink(input string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(input))
	if err != nil {
		return "", err
	}
	for _, key := range []string{"url", "URL", "target", "redirect"} {
		if value := parsed.Query().Get(key); value != "" {
			decoded, err := url.QueryUnescape(value)
			if err == nil {
				value = decoded
			}
			target, err := url.Parse(value)
			if err == nil && target.Scheme != "" && target.Host != "" {
				return target.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no complete destination URL found")
}
func MIME(input string) (string, string) {
	value := strings.ToLower(strings.TrimSpace(input))
	value = strings.TrimPrefix(value, ".")
	m := map[string]string{"html": "text/html", "css": "text/css", "js": "text/javascript", "json": "application/json", "xml": "application/xml", "txt": "text/plain", "pdf": "application/pdf", "png": "image/png", "jpg": "image/jpeg", "jpeg": "image/jpeg", "gif": "image/gif", "svg": "image/svg+xml", "webp": "image/webp", "mp3": "audio/mpeg", "mp4": "video/mp4", "zip": "application/zip", "csv": "text/csv"}
	if mime, ok := m[value]; ok {
		return value, mime
	}
	for ext, mime := range m {
		if value == mime {
			return ext, mime
		}
	}
	return "", ""
}
func UA(input string) []string {
	ua := strings.ToLower(input)
	browser := "Unknown"
	for _, pair := range [][2]string{{"edg/", "Microsoft Edge"}, {"firefox/", "Firefox"}, {"chrome/", "Chrome"}, {"safari/", "Safari"}} {
		if strings.Contains(ua, pair[0]) {
			browser = pair[1]
			break
		}
	}
	os := "Unknown"
	for _, pair := range [][2]string{{"windows", "Windows"}, {"android", "Android"}, {"iphone", "iOS"}, {"ipad", "iPadOS"}, {"mac os x", "macOS"}, {"linux", "Linux"}} {
		if strings.Contains(ua, pair[0]) {
			os = pair[1]
			break
		}
	}
	device := "Desktop"
	if strings.Contains(ua, "mobile") || strings.Contains(ua, "iphone") || strings.Contains(ua, "android") {
		device = "Mobile"
	}
	return []string{browser, os, device}
}
func OG(title, description, image, urlValue string) string {
	return fmt.Sprintf("<meta property=\"og:title\" content=\"%s\">\n<meta property=\"og:description\" content=\"%s\">\n<meta property=\"og:image\" content=\"%s\">\n<meta property=\"og:url\" content=\"%s\">\n<meta property=\"og:type\" content=\"website\">\n", html.EscapeString(title), html.EscapeString(description), html.EscapeString(image), html.EscapeString(urlValue))
}
func Status(code string) (string, string) {
	n, err := strconv.Atoi(code)
	if err != nil {
		return "", ""
	}
	m := map[int]string{100: "Continue", 200: "OK", 201: "Created", 204: "No Content", 301: "Moved Permanently", 302: "Found", 304: "Not Modified", 400: "Bad Request", 401: "Unauthorized", 403: "Forbidden", 404: "Not Found", 405: "Method Not Allowed", 408: "Request Timeout", 409: "Conflict", 418: "I'm a teapot", 429: "Too Many Requests", 500: "Internal Server Error", 502: "Bad Gateway", 503: "Service Unavailable", 504: "Gateway Timeout"}
	return code, m[n]
}
