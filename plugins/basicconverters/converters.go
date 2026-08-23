// Package basicconverters implements bounded, local-only conversion helpers.
package basicconverters

import (
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

// IntegerValues returns one integer rendered in the common positional bases.
func IntegerValues(input string, base int) (map[string]string, error) {
	if base != 2 && base != 8 && base != 10 && base != 16 {
		return nil, fmt.Errorf("unsupported input base")
	}
	input = strings.TrimSpace(input)
	if input == "" || len(input) > 1024 {
		return nil, fmt.Errorf("integer input must contain between 1 and 1024 characters")
	}
	value, ok := new(big.Int).SetString(input, base)
	if !ok {
		return nil, fmt.Errorf("input is not a valid base-%d integer", base)
	}
	return map[string]string{
		"binary":      value.Text(2),
		"octal":       value.Text(8),
		"decimal":     value.Text(10),
		"hexadecimal": strings.ToUpper(value.Text(16)),
	}, nil
}

var romanValues = []struct {
	value  int
	symbol string
}{
	{1000, "M"}, {900, "CM"}, {500, "D"}, {400, "CD"}, {100, "C"}, {90, "XC"},
	{50, "L"}, {40, "XL"}, {10, "X"}, {9, "IX"}, {5, "V"}, {4, "IV"}, {1, "I"},
}

// ToRoman converts a standard positive integer in the inclusive range 1-3999.
func ToRoman(value int) (string, error) {
	if value < 1 || value > 3999 {
		return "", fmt.Errorf("Roman numerals support whole numbers from 1 through 3999")
	}
	var output strings.Builder
	for _, part := range romanValues {
		for value >= part.value {
			output.WriteString(part.symbol)
			value -= part.value
		}
	}
	return output.String(), nil
}

// FromRoman parses a canonical standard Roman numeral.
func FromRoman(input string) (int, error) {
	input = strings.ToUpper(strings.TrimSpace(input))
	if input == "" || len(input) > 15 {
		return 0, fmt.Errorf("enter one Roman numeral from I through MMMCMXCIX")
	}
	remaining, total := input, 0
	for _, part := range romanValues {
		for strings.HasPrefix(remaining, part.symbol) {
			total += part.value
			remaining = strings.TrimPrefix(remaining, part.symbol)
		}
	}
	if remaining != "" || total < 1 || total > 3999 {
		return 0, fmt.Errorf("input is not a valid Roman numeral")
	}
	canonical, _ := ToRoman(total)
	if input != canonical {
		return 0, fmt.Errorf("input is not a canonical Roman numeral")
	}
	return total, nil
}

// Color is an sRGB color with 8-bit channels.
type Color struct{ R, G, B int }

var rgbPattern = regexp.MustCompile(`(?i)^rgb\(\s*(\d{1,3})\s*,\s*(\d{1,3})\s*,\s*(\d{1,3})\s*\)$`)
var hslPattern = regexp.MustCompile(`(?i)^hsl\(\s*(-?(?:\d+(?:\.\d*)?|\.\d+))\s*,\s*(\d{1,3}(?:\.\d*)?|\.\d+)\s*%\s*,\s*(\d{1,3}(?:\.\d*)?|\.\d+)\s*%\s*\)$`)

// ParseColor accepts #RGB, #RRGGBB, rgb(r,g,b), and hsl(h,s%,l%).
func ParseColor(input string) (Color, error) {
	input = strings.TrimSpace(input)
	if strings.HasPrefix(input, "#") {
		hex := strings.TrimPrefix(input, "#")
		if len(hex) == 3 {
			hex = strings.Repeat(string(hex[0]), 2) + strings.Repeat(string(hex[1]), 2) + strings.Repeat(string(hex[2]), 2)
		}
		if len(hex) != 6 {
			return Color{}, fmt.Errorf("hex colors must use #RGB or #RRGGBB")
		}
		value, err := strconv.ParseUint(hex, 16, 32)
		if err != nil {
			return Color{}, fmt.Errorf("hex color contains invalid characters")
		}
		return Color{R: int(value >> 16), G: int(value >> 8 & 0xff), B: int(value & 0xff)}, nil
	}
	if match := rgbPattern.FindStringSubmatch(input); match != nil {
		channels := [3]int{}
		for index := range channels {
			value, _ := strconv.Atoi(match[index+1])
			if value > 255 {
				return Color{}, fmt.Errorf("RGB channels must be between 0 and 255")
			}
			channels[index] = value
		}
		return Color{R: channels[0], G: channels[1], B: channels[2]}, nil
	}
	if match := hslPattern.FindStringSubmatch(input); match != nil {
		hue, _ := strconv.ParseFloat(match[1], 64)
		saturation, _ := strconv.ParseFloat(match[2], 64)
		lightness, _ := strconv.ParseFloat(match[3], 64)
		if saturation > 100 || lightness > 100 {
			return Color{}, fmt.Errorf("HSL saturation and lightness must be between 0%% and 100%%")
		}
		return hslToRGB(hue, saturation/100, lightness/100), nil
	}
	return Color{}, fmt.Errorf("enter #RGB, #RRGGBB, rgb(r, g, b), or hsl(h, s%%, l%%)")
}

func (c Color) Hex() string { return fmt.Sprintf("#%02X%02X%02X", c.R, c.G, c.B) }
func (c Color) RGB() string { return fmt.Sprintf("rgb(%d, %d, %d)", c.R, c.G, c.B) }

func (c Color) HSL() string {
	r, g, b := float64(c.R)/255, float64(c.G)/255, float64(c.B)/255
	maximum, minimum := math.Max(r, math.Max(g, b)), math.Min(r, math.Min(g, b))
	delta, lightness := maximum-minimum, (maximum+minimum)/2
	if delta == 0 {
		return fmt.Sprintf("hsl(0, 0%%, %.1f%%)", lightness*100)
	}
	saturation := delta / (1 - math.Abs(2*lightness-1))
	hue := 0.0
	switch maximum {
	case r:
		hue = math.Mod((g-b)/delta, 6)
	case g:
		hue = (b-r)/delta + 2
	default:
		hue = (r-g)/delta + 4
	}
	hue *= 60
	if hue < 0 {
		hue += 360
	}
	return fmt.Sprintf("hsl(%s, %s%%, %s%%)", decimal(hue), decimal(saturation*100), decimal(lightness*100))
}

func hslToRGB(hue, saturation, lightness float64) Color {
	hue = math.Mod(hue, 360)
	if hue < 0 {
		hue += 360
	}
	chroma := (1 - math.Abs(2*lightness-1)) * saturation
	x := chroma * (1 - math.Abs(math.Mod(hue/60, 2)-1))
	match := lightness - chroma/2
	r, g, b := 0.0, 0.0, 0.0
	switch {
	case hue < 60:
		r, g = chroma, x
	case hue < 120:
		r, g = x, chroma
	case hue < 180:
		g, b = chroma, x
	case hue < 240:
		g, b = x, chroma
	case hue < 300:
		r, b = x, chroma
	default:
		r, b = chroma, x
	}
	return Color{R: int(math.Round((r + match) * 255)), G: int(math.Round((g + match) * 255)), B: int(math.Round((b + match) * 255))}
}

func decimal(value float64) string {
	return strings.TrimRight(strings.TrimRight(strconv.FormatFloat(value, 'f', 1, 64), "0"), ".")
}
