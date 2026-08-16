// Package unicodeutils provides local Unicode escape transformations.
package unicodeutils

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

// Encode escapes every Unicode code point as a JavaScript-style Unicode escape.
// Non-BMP values are represented as their UTF-16 surrogate pair.
func Encode(input string) string {
	var output strings.Builder
	for _, value := range input {
		if value <= 0xFFFF {
			fmt.Fprintf(&output, "\\u%04X", value)
			continue
		}
		for _, unit := range utf16.Encode([]rune{value}) {
			fmt.Fprintf(&output, "\\u%04X", unit)
		}
	}
	return output.String()
}

// Decode replaces valid \uXXXX and \UXXXXXXXX escape sequences. Other text,
// including ordinary backslashes, remains unchanged.
func Decode(input string) (string, error) {
	var output strings.Builder
	for index := 0; index < len(input); {
		if input[index] != '\\' || index+1 >= len(input) {
			output.WriteByte(input[index])
			index++
			continue
		}
		switch input[index+1] {
		case 'u':
			if index+6 > len(input) {
				return "", fmt.Errorf("incomplete \\u escape at character %d", index+1)
			}
			unit, err := parseHex(input[index+2 : index+6])
			if err != nil {
				return "", fmt.Errorf("invalid \\u escape at character %d", index+1)
			}
			if utf16.IsSurrogate(rune(unit)) {
				if index+12 > len(input) || input[index+6:index+8] != "\\u" {
					return "", fmt.Errorf("unpaired surrogate at character %d", index+1)
				}
				second, err := parseHex(input[index+8 : index+12])
				if err != nil {
					return "", fmt.Errorf("invalid surrogate pair at character %d", index+1)
				}
				value := utf16.DecodeRune(rune(unit), rune(second))
				if value == utf8.RuneError {
					return "", fmt.Errorf("invalid surrogate pair at character %d", index+1)
				}
				output.WriteRune(value)
				index += 12
				continue
			}
			output.WriteRune(rune(unit))
			index += 6
		case 'U':
			if index+10 > len(input) {
				return "", fmt.Errorf("incomplete \\U escape at character %d", index+1)
			}
			value, err := parseHex(input[index+2 : index+10])
			if err != nil || value > utf8.MaxRune || (value >= 0xD800 && value <= 0xDFFF) {
				return "", fmt.Errorf("invalid \\U escape at character %d", index+1)
			}
			output.WriteRune(rune(value))
			index += 10
		default:
			output.WriteByte(input[index])
			index++
		}
	}
	return output.String(), nil
}

func parseHex(value string) (uint64, error) {
	return strconv.ParseUint(value, 16, 32)
}
