// Package qrcode generates local QR-code PNG images.
package qrcode

import (
	"fmt"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

// Generate creates a PNG QR code at size pixels with the requested recovery
// level. The data is encoded locally and not sent anywhere.
func Generate(data, recovery string, size int) ([]byte, error) {
	if strings.TrimSpace(data) == "" {
		return nil, fmt.Errorf("text or URL is required")
	}
	if size < 64 || size > 1024 {
		return nil, fmt.Errorf("image size must be between 64 and 1024 pixels")
	}
	levels := map[string]qrcode.RecoveryLevel{"low": qrcode.Low, "medium": qrcode.Medium, "high": qrcode.Highest, "highest": qrcode.Highest}
	level, ok := levels[recovery]
	if !ok {
		return nil, fmt.Errorf("unsupported error recovery level")
	}
	encoded, err := qrcode.Encode(data, level, size)
	if err != nil {
		return nil, fmt.Errorf("generate QR code: %w", err)
	}
	return encoded, nil
}
