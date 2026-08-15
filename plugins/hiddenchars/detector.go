// Package hiddenchars finds potentially invisible or confusing Unicode code points.
package hiddenchars

import (
	"fmt"
	"unicode"
)

type Finding struct {
	Line        int
	Column      int
	CodePoint   string
	Description string
}

// Detect reports formatting characters, non-standard spaces, and control
// characters other than normal tabs and line breaks. Line and column values
// are one-based Unicode character positions.
func Detect(input string) []Finding {
	findings := make([]Finding, 0)
	line, column := 1, 1
	for _, character := range input {
		if suspicious(character) {
			findings = append(findings, Finding{Line: line, Column: column, CodePoint: fmt.Sprintf("U+%04X", character), Description: description(character)})
		}
		if character == '\n' {
			line++
			column = 1
		} else {
			column++
		}
	}
	return findings
}

func suspicious(character rune) bool {
	if character == '\n' || character == '\r' || character == '\t' || character == ' ' {
		return false
	}
	return unicode.Is(unicode.Cf, character) || unicode.IsControl(character) || unicode.Is(unicode.Zs, character) || character == 0x00AD
}

func description(character rune) string {
	if known, ok := knownDescriptions[character]; ok {
		return known
	}
	if unicode.Is(unicode.Cf, character) {
		return "Unicode format character"
	}
	if unicode.Is(unicode.Zs, character) {
		return "Non-standard space separator"
	}
	return "Control character"
}

var knownDescriptions = map[rune]string{
	0x00A0: "Non-breaking space",
	0x00AD: "Soft hyphen",
	0x034F: "Combining grapheme joiner",
	0x061C: "Arabic letter mark",
	0x180E: "Mongolian vowel separator",
	0x200B: "Zero-width space",
	0x200C: "Zero-width non-joiner",
	0x200D: "Zero-width joiner",
	0x200E: "Left-to-right mark",
	0x200F: "Right-to-left mark",
	0x202A: "Left-to-right embedding",
	0x202B: "Right-to-left embedding",
	0x202C: "Pop directional formatting",
	0x202D: "Left-to-right override",
	0x202E: "Right-to-left override",
	0x202F: "Narrow non-breaking space",
	0x2060: "Word joiner",
	0x2066: "Left-to-right isolate",
	0x2067: "Right-to-left isolate",
	0x2068: "First strong isolate",
	0x2069: "Pop directional isolate",
	0x206A: "Inhibit symmetric swapping",
	0x206B: "Activate symmetric swapping",
	0x206C: "Inhibit Arabic form shaping",
	0x206D: "Activate Arabic form shaping",
	0x206E: "National digit shapes",
	0x206F: "Nominal digit shapes",
	0xFEFF: "Zero-width no-break space / byte order mark",
}
