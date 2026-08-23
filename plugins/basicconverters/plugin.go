package basicconverters

import (
	"fmt"
	"strconv"
	"strings"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
)

func RunIntegerBaseConverter() {
	cmdry.Run(plugin("integer-base-converter", "Integer Base Converter", "Convert pasted integers among binary, octal, decimal, and hexadecimal locally.", []string{"integer", "binary", "octal", "decimal", "hexadecimal", "base", "number"}, integerForm, integerConvert))
}

func RunRomanNumeralConverter() {
	cmdry.Run(plugin("roman-numeral-converter", "Roman Numeral Converter", "Convert values between standard Roman numerals and decimal numbers locally.", []string{"roman", "numeral", "decimal", "number", "convert"}, romanForm, romanConvert))
}

func RunColorConverter() {
	cmdry.Run(plugin("color-converter", "Color Converter", "Convert hex, RGB, and HSL color values locally.", []string{"color", "hex", "rgb", "hsl", "hsv", "convert"}, colorForm, colorConvert))
}

func plugin(id, name, description string, terms []string, overview, action cmdry.Handler) cmdry.Plugin {
	return cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey." + id, Name: name, Version: "0.1.0", Description: description, Category: "converter", Icon: "convert", SearchTerms: terms, Pages: []cmdry.Page{{ID: "overview", Name: "Convert", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New conversion", Method: "read"}, {ID: "convert", Name: "Convert", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "convert": action}}
}

func integerForm(cmdry.Request) (cmdry.View, error) {
	return form("Integer Base Converter", "Convert one whole number. The input has no 0x, 0o, or 0b prefix; choose its base below.", "Convert integer", []cmdry.Field{{Name: "input", Label: "Integer", Type: "text", Value: "255", Required: true}, {Name: "base", Label: "Input base", Type: "select", Value: "10", Options: []cmdry.Option{{Value: "2", Label: "Binary (base 2)"}, {Value: "8", Label: "Octal (base 8)"}, {Value: "10", Label: "Decimal (base 10)"}, {Value: "16", Label: "Hexadecimal (base 16)"}}}}), nil
}

func integerConvert(request cmdry.Request) (cmdry.View, error) {
	base, err := strconv.Atoi(fmt.Sprint(request.Params["base"]))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("choose a supported input base")
	}
	values, err := IntegerValues(fmt.Sprint(request.Params["input"]), base)
	if err != nil {
		return cmdry.View{}, err
	}
	rows := []map[string]any{{"base": "Binary (base 2)", "value": values["binary"]}, {"base": "Octal (base 8)", "value": values["octal"]}, {"base": "Decimal (base 10)", "value": values["decimal"]}, {"base": "Hexadecimal (base 16)", "value": values["hexadecimal"]}}
	return result("Integer conversion", "Equivalent values are shown in common integer bases.", "integer-base", rows), nil
}

func romanForm(cmdry.Request) (cmdry.View, error) {
	return form("Roman Numeral Converter", "Use standard Roman numerals from I through MMMCMXCIX, or a whole decimal number from 1 through 3999.", "Convert numeral", []cmdry.Field{{Name: "operation", Label: "Convert", Type: "select", Value: "decimal-to-roman", Options: []cmdry.Option{{Value: "decimal-to-roman", Label: "Decimal to Roman"}, {Value: "roman-to-decimal", Label: "Roman to decimal"}}}, {Name: "input", Label: "Number or numeral", Type: "text", Value: "2026", Required: true}}), nil
}

func romanConvert(request cmdry.Request) (cmdry.View, error) {
	input, operation := strings.TrimSpace(fmt.Sprint(request.Params["input"])), fmt.Sprint(request.Params["operation"])
	var decimalValue int
	var roman string
	var err error
	switch operation {
	case "decimal-to-roman":
		decimalValue, err = strconv.Atoi(input)
		if err == nil {
			roman, err = ToRoman(decimalValue)
		}
	case "roman-to-decimal":
		decimalValue, err = FromRoman(input)
		if err == nil {
			roman, _ = ToRoman(decimalValue)
		}
	default:
		return cmdry.View{}, fmt.Errorf("choose a supported conversion")
	}
	if err != nil {
		return cmdry.View{}, err
	}
	return result("Roman numeral conversion", "The Roman representation is canonical standard notation.", "roman-numeral", []map[string]any{{"base": "Decimal", "value": strconv.Itoa(decimalValue)}, {"base": "Roman numeral", "value": roman}}), nil
}

func colorForm(cmdry.Request) (cmdry.View, error) {
	return form("Color Converter", "Enter #RGB, #RRGGBB, rgb(r, g, b), or hsl(h, s%, l%).", "Convert color", []cmdry.Field{{Name: "input", Label: "Color", Type: "text", Value: "#4AB7AA", Required: true}}), nil
}

func colorConvert(request cmdry.Request) (cmdry.View, error) {
	color, err := ParseColor(fmt.Sprint(request.Params["input"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return result("Color conversion", "Equivalent color values are shown in common web formats.", "color", []map[string]any{{"base": "Hex", "value": color.Hex()}, {"base": "RGB", "value": color.RGB()}, {"base": "HSL", "value": color.HSL()}}), nil
}

func form(title, description, submit string, fields []cmdry.Field) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "form", Title: title, Action: "convert", Submit: submit, Description: description + " Your data never leaves this device.", Fields: fields}}}
}

func result(title, message, id string, rows []map[string]any) cmdry.View {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "Conversion complete", Message: message}, {Type: "table", ID: id, Columns: []cmdry.Column{{Key: "base", Label: "Format"}, {Key: "value", Label: "Value"}}, Rows: rows}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Convert another value"}}}}}
}
