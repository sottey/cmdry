package texttools

import (
	"fmt"
	"strings"
	"unicode"
)

// RotateText shifts Unicode characters by positions. Positive positions move
// the first characters to the end; negative positions move the last forward.
func RotateText(input string, positions int) string {
	chars := []rune(input)
	if len(chars) == 0 {
		return ""
	}
	shift := positions % len(chars)
	if shift < 0 {
		shift += len(chars)
	}
	return string(append(chars[shift:], chars[:shift]...))
}

var morseAlphabet = map[rune]string{
	'a': ".-", 'b': "-...", 'c': "-.-.", 'd': "-..", 'e': ".", 'f': "..-.", 'g': "--.", 'h': "....", 'i': "..", 'j': ".---", 'k': "-.-", 'l': ".-..", 'm': "--", 'n': "-.", 'o': "---", 'p': ".--.", 'q': "--.-", 'r': ".-.", 's': "...", 't': "-", 'u': "..-", 'v': "...-", 'w': ".--", 'x': "-..-", 'y': "-.--", 'z': "--..",
	'0': "-----", '1': ".----", '2': "..---", '3': "...--", '4': "....-", '5': ".....", '6': "-....", '7': "--...", '8': "---..", '9': "----.",
	'.': ".-.-.-", ',': "--..--", '?': "..--..", '!': "-.-.--", '/': "-..-.", '(': "-.--.", ')': "-.--.-", '&': ".-...", ':': "---...", ';': "-.-.-.", '=': "-...-", '+': ".-.-.", '-': "-....-", '_': "..--.-", '"': ".-..-.", '$': "...-..-", '@': ".--.-.",
}

// ToMorse encodes supported letters, digits, and punctuation. Word breaks are
// represented as a slash and unsupported characters are rejected explicitly.
func ToMorse(input string) (string, error) {
	words := strings.Fields(input)
	if len(words) == 0 {
		return "", fmt.Errorf("enter text to encode")
	}
	encoded := make([]string, 0, len(words))
	for _, word := range words {
		letters := make([]string, 0, len(word))
		for _, char := range strings.ToLower(word) {
			code, ok := morseAlphabet[char]
			if !ok {
				return "", fmt.Errorf("unsupported character %q", char)
			}
			letters = append(letters, code)
		}
		encoded = append(encoded, strings.Join(letters, " "))
	}
	return strings.Join(encoded, " / "), nil
}

// CensorText replaces case-insensitive whole-word matches with replacement.
func CensorText(input, words, replacement string) (string, int, error) {
	terms := strings.FieldsFunc(words, func(char rune) bool { return char == ',' || char == '\n' || char == '\r' })
	lookup := make(map[string]struct{}, len(terms))
	for _, term := range terms {
		term = strings.TrimSpace(term)
		if term != "" {
			lookup[strings.ToLower(term)] = struct{}{}
		}
	}
	if len(lookup) == 0 {
		return "", 0, fmt.Errorf("enter at least one word to censor")
	}
	if replacement == "" {
		replacement = "***"
	}
	var output, word strings.Builder
	count := 0
	flush := func() {
		if word.Len() == 0 {
			return
		}
		value := word.String()
		if _, ok := lookup[strings.ToLower(value)]; ok {
			output.WriteString(replacement)
			count++
		} else {
			output.WriteString(value)
		}
		word.Reset()
	}
	for _, char := range input {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			word.WriteRune(char)
		} else {
			flush()
			output.WriteRune(char)
		}
	}
	flush()
	return output.String(), count, nil
}

// QuoteText surrounds each non-blank input line with one selected quote style.
func QuoteText(input, style string) (string, int, error) {
	left, right := "\"", "\""
	switch style {
	case "single":
		left, right = "'", "'"
	case "double":
	case "curly":
		left, right = "“", "”"
	default:
		return "", 0, fmt.Errorf("choose a valid quote style")
	}
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	output := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			output = append(output, left+line+right)
		}
	}
	if len(output) == 0 {
		return "", 0, fmt.Errorf("enter at least one non-blank line")
	}
	return strings.Join(output, "\n"), len(output), nil
}
