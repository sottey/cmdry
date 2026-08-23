// Package imagepriority2 provides local QR and SVG image generators.
package imagepriority2

import (
	"fmt"
	"regexp"
	"strings"
)

func WiFiPayload(ssid, password, security string, hidden bool) (string, error) {
	ssid = strings.TrimSpace(ssid)
	if ssid == "" || len(ssid) > 128 {
		return "", fmt.Errorf("network name must contain 1 to 128 characters")
	}
	if len(password) > 256 {
		return "", fmt.Errorf("password must not exceed 256 characters")
	}
	switch security {
	case "WPA", "WEP":
		if password == "" {
			return "", fmt.Errorf("a password is required for %s security", security)
		}
	case "nopass":
		if password != "" {
			return "", fmt.Errorf("open networks must not include a password")
		}
	default:
		return "", fmt.Errorf("choose a supported security type")
	}
	payload := "WIFI:T:" + security + ";S:" + wifiEscape(ssid) + ";"
	if password != "" {
		payload += "P:" + wifiEscape(password) + ";"
	}
	if hidden {
		payload += "H:true;"
	}
	return payload + ";", nil
}
func wifiEscape(value string) string {
	replacer := strings.NewReplacer("\\", "\\\\", ";", "\\;", ",", "\\,", ":", "\\:")
	return replacer.Replace(value)
}

func SVGPlaceholder(width, height int, label, background, foreground string) ([]byte, error) {
	if width < 1 || width > 4096 || height < 1 || height > 4096 {
		return nil, fmt.Errorf("dimensions must be between 1 and 4,096 pixels")
	}
	label = strings.TrimSpace(label)
	if len(label) > 120 {
		return nil, fmt.Errorf("label must not exceed 120 characters")
	}
	if label == "" {
		label = fmt.Sprintf("%d × %d", width, height)
	}
	if !validColor(background) || !validColor(foreground) {
		return nil, fmt.Errorf("colors must be 3- or 6-digit hex values, such as #0f172a")
	}
	fontSize := width / 12
	if height/4 < fontSize {
		fontSize = height / 4
	}
	if fontSize < 12 {
		fontSize = 12
	}
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d" role="img" aria-label="%s"><rect width="100%%" height="100%%" fill="%s"/><text x="50%%" y="50%%" fill="%s" dominant-baseline="middle" text-anchor="middle" font-family="system-ui, sans-serif" font-size="%d">%s</text></svg>`+"\n", width, height, width, height, xmlEscape(label), background, foreground, fontSize, xmlEscape(label))
	return []byte(svg), nil
}

var hexColor = regexp.MustCompile(`^#[0-9a-fA-F]{3}([0-9a-fA-F]{3})?$`)

func validColor(value string) bool { return hexColor.MatchString(strings.TrimSpace(value)) }
func xmlEscape(value string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "\"", "&quot;", "'", "&apos;").Replace(value)
}
