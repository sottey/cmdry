// Package textpriority2 implements bounded, local-only text utilities.
package textpriority2

import (
	"fmt"
	"strings"
	"unicode"
)

var natoEncode = map[rune]string{
	'A': "Alfa", 'B': "Bravo", 'C': "Charlie", 'D': "Delta", 'E': "Echo", 'F': "Foxtrot", 'G': "Golf", 'H': "Hotel", 'I': "India", 'J': "Juliett", 'K': "Kilo", 'L': "Lima", 'M': "Mike", 'N': "November", 'O': "Oscar", 'P': "Papa", 'Q': "Quebec", 'R': "Romeo", 'S': "Sierra", 'T': "Tango", 'U': "Uniform", 'V': "Victor", 'W': "Whiskey", 'X': "X-ray", 'Y': "Yankee", 'Z': "Zulu",
	'0': "Zero", '1': "One", '2': "Two", '3': "Three", '4': "Four", '5': "Five", '6': "Six", '7': "Seven", '8': "Eight", '9': "Nine",
}

var natoDecode = func() map[string]rune {
	values := make(map[string]rune, len(natoEncode))
	for letter, word := range natoEncode {
		values[strings.ToUpper(word)] = letter
	}
	return values
}()

// NATOEncode writes letters and digits as NATO phonetic words.
func NATOEncode(input string) (string, error) {
	if len([]rune(input)) == 0 || len([]rune(input)) > 1000 {
		return "", fmt.Errorf("enter between 1 and 1,000 characters")
	}
	words := make([]string, 0, len(input))
	for _, character := range strings.ToUpper(input) {
		if unicode.IsSpace(character) {
			if len(words) > 0 && words[len(words)-1] != "/" {
				words = append(words, "/")
			}
			continue
		}
		word, ok := natoEncode[character]
		if !ok {
			return "", fmt.Errorf("NATO encoding supports letters, digits, and spaces")
		}
		words = append(words, word)
	}
	return strings.Trim(strings.Join(words, " "), " /"), nil
}

// NATODecode parses whitespace-separated NATO words; slash separates words.
func NATODecode(input string) (string, error) {
	tokens := strings.Fields(strings.TrimSpace(input))
	if len(tokens) == 0 || len(tokens) > 1000 {
		return "", fmt.Errorf("enter between 1 and 1,000 NATO words")
	}
	var output strings.Builder
	for _, token := range tokens {
		if token == "/" {
			if output.Len() > 0 && !strings.HasSuffix(output.String(), " ") {
				output.WriteByte(' ')
			}
			continue
		}
		letter, ok := natoDecode[strings.ToUpper(token)]
		if !ok {
			return "", fmt.Errorf("%q is not a recognized NATO word", token)
		}
		output.WriteRune(letter)
	}
	return strings.TrimSpace(output.String()), nil
}

// ASCIIBinaryEncode converts ASCII text to spaced eight-bit binary bytes.
func ASCIIBinaryEncode(input string) (string, error) {
	if input == "" || len(input) > 4096 {
		return "", fmt.Errorf("enter between 1 and 4,096 ASCII characters")
	}
	values := make([]string, 0, len(input))
	for _, character := range input {
		if character > 127 {
			return "", fmt.Errorf("ASCII binary supports only ASCII characters")
		}
		values = append(values, fmt.Sprintf("%08b", character))
	}
	return strings.Join(values, " "), nil
}

// ASCIIBinaryDecode parses whitespace-separated eight-bit ASCII values.
func ASCIIBinaryDecode(input string) (string, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 || len(parts) > 4096 {
		return "", fmt.Errorf("enter between 1 and 4,096 eight-bit values")
	}
	var output strings.Builder
	for _, part := range parts {
		if len(part) != 8 {
			return "", fmt.Errorf("each binary value must contain exactly eight bits")
		}
		value := byte(0)
		for _, bit := range part {
			if bit != '0' && bit != '1' {
				return "", fmt.Errorf("binary values may contain only 0 and 1")
			}
			value = value<<1 | byte(bit-'0')
		}
		if value > 127 {
			return "", fmt.Errorf("binary values must represent ASCII bytes from 0 through 127")
		}
		output.WriteByte(value)
	}
	return output.String(), nil
}

var loremWords = []string{"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit", "sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore", "magna", "aliqua"}

// LoremIpsum produces a deterministic number of conventional placeholder paragraphs.
func LoremIpsum(paragraphs int) (string, error) {
	if paragraphs < 1 || paragraphs > 20 {
		return "", fmt.Errorf("paragraph count must be between 1 and 20")
	}
	output := make([]string, paragraphs)
	for paragraph := range output {
		words := make([]string, 48)
		for index := range words {
			words[index] = loremWords[(paragraph*7+index)%len(loremWords)]
		}
		words[0] = strings.ToUpper(words[0][:1]) + words[0][1:]
		output[paragraph] = strings.Join(words[:16], " ") + ". " + strings.Join(words[16:32], " ") + ". " + strings.Join(words[32:], " ") + "."
	}
	return strings.Join(output, "\n\n"), nil
}

// ROT47 applies the same reversible ROT47 transform in either direction.
func ROT47(input string) (string, error) {
	if len(input) == 0 || len(input) > 10000 {
		return "", fmt.Errorf("enter between 1 and 10,000 characters")
	}
	var output strings.Builder
	for _, character := range input {
		if character >= 33 && character <= 126 {
			output.WriteRune(33 + (character-33+47)%94)
		} else {
			output.WriteRune(character)
		}
	}
	return output.String(), nil
}

// Numeronymize shortens words longer than three letters, retaining punctuation and spacing.
func Numeronymize(input string) (string, error) {
	if len([]rune(input)) == 0 || len([]rune(input)) > 10000 {
		return "", fmt.Errorf("enter between 1 and 10,000 characters")
	}
	var output strings.Builder
	for start := 0; start < len(input); {
		if !isASCIIAlpha(input[start]) {
			output.WriteByte(input[start])
			start++
			continue
		}
		end := start
		for end < len(input) && isASCIIAlpha(input[end]) {
			end++
		}
		word := input[start:end]
		if len(word) > 3 {
			output.WriteByte(word[0])
			output.WriteString(fmt.Sprint(len(word) - 2))
			output.WriteByte(word[len(word)-1])
		} else {
			output.WriteString(word)
		}
		start = end
	}
	return output.String(), nil
}

func isASCIIAlpha(value byte) bool {
	return value >= 'A' && value <= 'Z' || value >= 'a' && value <= 'z'
}

var asciiGlyphs = map[rune][]string{
	'A': {" # ", "# #", "###", "# #", "# #"}, 'B': {"## ", "# #", "## ", "# #", "## "}, 'C': {" ##", "#  ", "#  ", "#  ", " ##"}, 'D': {"## ", "# #", "# #", "# #", "## "}, 'E': {"###", "#  ", "## ", "#  ", "###"}, 'F': {"###", "#  ", "## ", "#  ", "#  "}, 'G': {" ##", "#  ", "# #", "# #", " ##"}, 'H': {"# #", "# #", "###", "# #", "# #"}, 'I': {"###", " # ", " # ", " # ", "###"}, 'J': {"  #", "  #", "  #", "# #", " # "}, 'K': {"# #", "## ", "#  ", "## ", "# #"}, 'L': {"#  ", "#  ", "#  ", "#  ", "###"}, 'M': {"# #", "###", "###", "# #", "# #"}, 'N': {"# #", "###", "###", "###", "# #"}, 'O': {" # ", "# #", "# #", "# #", " # "}, 'P': {"## ", "# #", "## ", "#  ", "#  "}, 'Q': {" # ", "# #", "# #", " ##", "  #"}, 'R': {"## ", "# #", "## ", "# #", "# #"}, 'S': {" ##", "#  ", " # ", "  #", "## "}, 'T': {"###", " # ", " # ", " # ", " # "}, 'U': {"# #", "# #", "# #", "# #", "###"}, 'V': {"# #", "# #", "# #", "# #", " # "}, 'W': {"# #", "# #", "###", "###", "# #"}, 'X': {"# #", "# #", " # ", "# #", "# #"}, 'Y': {"# #", "# #", " # ", " # ", " # "}, 'Z': {"###", "  #", " # ", "#  ", "###"},
	'0': {"###", "# #", "# #", "# #", "###"}, '1': {" # ", "## ", " # ", " # ", "###"}, '2': {"## ", "  #", " # ", "#  ", "###"}, '3': {"## ", "  #", " ##", "  #", "## "}, '4': {"# #", "# #", "###", "  #", "  #"}, '5': {"###", "#  ", "## ", "  #", "## "}, '6': {" ##", "#  ", "## ", "# #", " ##"}, '7': {"###", "  #", " # ", "#  ", "#  "}, '8': {" ##", "# #", " ##", "# #", " ##"}, '9': {" ##", "# #", " ###", "  #", "## "},
	' ': {"   ", "   ", "   ", "   ", "   "}, '-': {"   ", "   ", "###", "   ", "   "}, '!': {" # ", " # ", " # ", "   ", " # "}, '?': {"## ", "  #", " # ", "   ", " # "},
}

// ASCIIArt renders a bounded line of text with Cmdry's bundled block font.
func ASCIIArt(input string) (string, error) {
	letters := []rune(strings.ToUpper(strings.TrimSpace(input)))
	if len(letters) == 0 || len(letters) > 40 {
		return "", fmt.Errorf("enter between 1 and 40 characters")
	}
	lines := make([]string, 5)
	for _, letter := range letters {
		glyph, ok := asciiGlyphs[letter]
		if !ok {
			glyph = asciiGlyphs['?']
		}
		for row := range lines {
			lines[row] += glyph[row] + " "
		}
	}
	return strings.TrimRight(strings.Join(lines, "\n"), " "), nil
}

var emojiCatalog = []struct{ Value, Label string }{
	{"😀", "Smile — 😀"}, {"😂", "Laughing tears — 😂"}, {"😍", "Heart eyes — 😍"}, {"🤔", "Thinking — 🤔"}, {"😭", "Crying — 😭"}, {"👍", "Thumbs up — 👍"}, {"👎", "Thumbs down — 👎"}, {"🙏", "Thank you — 🙏"}, {"👏", "Clap — 👏"}, {"🎉", "Party — 🎉"}, {"🔥", "Fire — 🔥"}, {"✨", "Sparkles — ✨"}, {"✅", "Check mark — ✅"}, {"❌", "Cross mark — ❌"}, {"⚠️", "Warning — ⚠️"}, {"💡", "Idea — 💡"}, {"🚀", "Rocket — 🚀"}, {"💻", "Laptop — 💻"}, {"📌", "Pin — 📌"}, {"📦", "Package — 📦"}, {"🔒", "Lock — 🔒"}, {"🔑", "Key — 🔑"}, {"🌐", "Globe — 🌐"}, {"❤️", "Heart — ❤️"}, {"⭐", "Star — ⭐"}, {"🍕", "Pizza — 🍕"}, {"☕", "Coffee — ☕"}, {"🐶", "Dog — 🐶"}, {"🐱", "Cat — 🐱"}, {"🌈", "Rainbow — 🌈"},
}

// EmojiOptions returns a copy of the bundled local emoji catalog.
func EmojiOptions() []struct{ Value, Label string } {
	return append([]struct{ Value, Label string }{}, emojiCatalog...)
}
