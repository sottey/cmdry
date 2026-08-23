// Package datapriority2 implements bounded, local-only data conversion and validation helpers.
package datapriority2

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"

	"github.com/sottey/cmdry/plugins/conversiontools"
	"go.yaml.in/yaml/v3"
)

const maxInputBytes = 1 << 20

func ConvertYAMLFamily(input, operation string) ([]byte, error) {
	if len(input) == 0 || len(input) > maxInputBytes {
		return nil, fmt.Errorf("input must contain between 1 and %d bytes", maxInputBytes)
	}
	switch operation {
	case "yaml-json":
		value, err := parseYAML(input)
		if err != nil {
			return nil, err
		}
		return jsonOutput(value)
	case "yaml-toml":
		value, err := parseYAML(input)
		if err != nil {
			return nil, err
		}
		jsonValue, err := jsonOutput(value)
		if err != nil {
			return nil, err
		}
		return conversiontools.JSONToTOML(string(jsonValue))
	case "json-yaml":
		value, err := parseJSON(input)
		if err != nil {
			return nil, err
		}
		return yamlOutput(value)
	case "toml-yaml":
		jsonValue, err := conversiontools.TOMLToJSON(input)
		if err != nil {
			return nil, err
		}
		value, err := parseJSON(string(jsonValue))
		if err != nil {
			return nil, err
		}
		return yamlOutput(value)
	default:
		return nil, fmt.Errorf("choose a supported conversion")
	}
}

func parseYAML(input string) (any, error) {
	var value any
	if err := yaml.Unmarshal([]byte(input), &value); err != nil {
		return nil, fmt.Errorf("parse YAML: %w", err)
	}
	return yamlToJSONValue(value)
}
func yamlToJSONValue(value any) (any, error) {
	switch typed := value.(type) {
	case map[string]any:
		out := make(map[string]any, len(typed))
		for key, item := range typed {
			converted, err := yamlToJSONValue(item)
			if err != nil {
				return nil, err
			}
			out[key] = converted
		}
		return out, nil
	case map[any]any:
		out := make(map[string]any, len(typed))
		for key, item := range typed {
			stringKey, ok := key.(string)
			if !ok {
				return nil, fmt.Errorf("YAML mapping keys must be strings for conversion")
			}
			converted, err := yamlToJSONValue(item)
			if err != nil {
				return nil, err
			}
			out[stringKey] = converted
		}
		return out, nil
	case []any:
		out := make([]any, len(typed))
		for index, item := range typed {
			converted, err := yamlToJSONValue(item)
			if err != nil {
				return nil, err
			}
			out[index] = converted
		}
		return out, nil
	default:
		return value, nil
	}
}
func parseJSON(input string) (any, error) {
	decoder := json.NewDecoder(strings.NewReader(input))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("input must contain one JSON value")
	}
	return value, nil
}
func jsonOutput(value any) ([]byte, error) {
	output, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode JSON: %w", err)
	}
	return append(output, '\n'), nil
}
func yamlOutput(value any) ([]byte, error) {
	output, err := yaml.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("encode YAML: %w", err)
	}
	return output, nil
}

// XMLToJSON converts a simple XML document to a structured JSON representation.
func XMLToJSON(input string) ([]byte, error) {
	if len(input) == 0 || len(input) > maxInputBytes {
		return nil, fmt.Errorf("XML must contain between 1 and %d bytes", maxInputBytes)
	}
	decoder := xml.NewDecoder(strings.NewReader(input))
	decoder.Strict = true
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			return nil, fmt.Errorf("XML contains no root element")
		}
		if err != nil {
			return nil, fmt.Errorf("parse XML: %w", err)
		}
		if start, ok := token.(xml.StartElement); ok {
			counter := 0
			value, err := decodeElement(decoder, start, 1, &counter)
			if err != nil {
				return nil, err
			}
			for {
				trailing, err := decoder.Token()
				if err == io.EOF {
					return jsonOutput(map[string]any{start.Name.Local: value})
				}
				if err != nil {
					return nil, fmt.Errorf("parse XML: %w", err)
				}
				if data, ok := trailing.(xml.CharData); ok && strings.TrimSpace(string(data)) == "" {
					continue
				}
				if _, ok := trailing.(xml.Comment); ok {
					continue
				}
				return nil, fmt.Errorf("XML must contain exactly one root element")
			}
		}
	}
}
func decodeElement(decoder *xml.Decoder, start xml.StartElement, depth int, counter *int) (any, error) {
	if depth > 64 {
		return nil, fmt.Errorf("XML nesting exceeds 64 levels")
	}
	*counter++
	if *counter > 10000 {
		return nil, fmt.Errorf("XML exceeds 10,000 elements")
	}
	value := map[string]any{}
	for _, attribute := range start.Attr {
		value["@"+attribute.Name.Local] = attribute.Value
	}
	var text strings.Builder
	for {
		token, err := decoder.Token()
		if err != nil {
			return nil, fmt.Errorf("parse XML: %w", err)
		}
		switch typed := token.(type) {
		case xml.StartElement:
			child, err := decodeElement(decoder, typed, depth+1, counter)
			if err != nil {
				return nil, err
			}
			if existing, found := value[typed.Name.Local]; found {
				switch children := existing.(type) {
				case []any:
					value[typed.Name.Local] = append(children, child)
				default:
					value[typed.Name.Local] = []any{children, child}
				}
			} else {
				value[typed.Name.Local] = child
			}
		case xml.CharData:
			text.Write([]byte(typed))
		case xml.EndElement:
			trimmed := strings.TrimSpace(text.String())
			if len(value) == 0 {
				return trimmed, nil
			}
			if trimmed != "" {
				value["#text"] = trimmed
			}
			return value, nil
		}
	}
}

type Phone struct{ Input, E164, National, International, Region, Note string }

// ParsePhone supports NANP formatting and structural E.164 normalization.
func ParsePhone(input, region string) (Phone, error) {
	input = strings.TrimSpace(input)
	if input == "" || len(input) > 64 {
		return Phone{}, fmt.Errorf("enter a phone number up to 64 characters")
	}
	digits, plus := strings.Builder{}, strings.HasPrefix(input, "+")
	for _, character := range input {
		if unicode.IsDigit(character) {
			digits.WriteRune(character)
		} else if !strings.ContainsRune("+()-. ", character) {
			return Phone{}, fmt.Errorf("phone number contains unsupported characters")
		}
	}
	value := digits.String()
	if region == "nanp" {
		if strings.HasPrefix(value, "1") && len(value) == 11 {
			value = value[1:]
		}
		if len(value) != 10 || value[0] < '2' || value[0] > '9' || value[3] < '2' || value[3] > '9' {
			return Phone{}, fmt.Errorf("NANP numbers must contain a valid 10-digit area code and exchange")
		}
		return Phone{Input: input, E164: "+1" + value, National: fmt.Sprintf("(%s) %s-%s", value[:3], value[3:6], value[6:]), International: "+1 " + value[:3] + " " + value[3:6] + " " + value[6:], Region: "US / Canada NANP", Note: "Structural validation only; this does not confirm a number is assigned."}, nil
	}
	if !plus || len(value) < 7 || len(value) > 15 || value[0] == '0' {
		return Phone{}, fmt.Errorf("international E.164 input must start with + and contain 7 to 15 digits")
	}
	return Phone{Input: input, E164: "+" + value, International: "+" + value, Region: "International E.164", Note: "Normalized structurally; country-specific assignment and formatting are not verified."}, nil
}

type IBAN struct {
	Normalized, CountryCode, Country, BBAN string
	Fields                                 []Field
}
type Field struct{ Label, Value string }

// ParseIBAN validates the ISO 13616 mod-97 checksum and shows known layout fields.
func ParseIBAN(input string) (IBAN, error) {
	normalized := strings.ToUpper(strings.Map(func(character rune) rune {
		if unicode.IsLetter(character) || unicode.IsDigit(character) {
			return character
		}
		return -1
	}, input))
	if len(normalized) < 15 || len(normalized) > 34 || !regexp.MustCompile(`^[A-Z]{2}[0-9]{2}[A-Z0-9]+$`).MatchString(normalized) {
		return IBAN{}, fmt.Errorf("enter an IBAN with 15 to 34 letters and digits")
	}
	if expectedLength, known := ibanLengths[normalized[:2]]; known && len(normalized) != expectedLength {
		return IBAN{}, fmt.Errorf("%s IBANs must contain %d characters", normalized[:2], expectedLength)
	}
	if ibanMod97(normalized[4:]+normalized[:4]) != 1 {
		return IBAN{}, fmt.Errorf("IBAN checksum is invalid")
	}
	result := IBAN{Normalized: normalized, CountryCode: normalized[:2], Country: ibanCountries[normalized[:2]], BBAN: normalized[4:]}
	if result.Country == "" {
		result.Country = "Country code recognized by the IBAN format; name not bundled"
	}
	result.Fields = ibanFields(normalized)
	return result, nil
}
func ibanMod97(value string) int {
	remainder := 0
	for _, character := range value {
		if character >= 'A' && character <= 'Z' {
			remainder = (remainder*100 + int(character-'A'+10)) % 97
		} else {
			remainder = (remainder*10 + int(character-'0')) % 97
		}
	}
	return remainder
}
func ibanFields(value string) []Field {
	switch value[:2] {
	case "DE":
		return []Field{{"Bank code", value[4:12]}, {"Account number", value[12:]}}
	case "GB":
		return []Field{{"Bank code", value[4:8]}, {"Sort code", value[8:14]}, {"Account number", value[14:]}}
	case "FR":
		return []Field{{"Bank code", value[4:9]}, {"Branch code", value[9:14]}, {"Account number", value[14:25]}, {"National check digits", value[25:]}}
	case "ES":
		return []Field{{"Bank code", value[4:8]}, {"Branch code", value[8:12]}, {"National check digits", value[12:14]}, {"Account number", value[14:]}}
	case "NL":
		return []Field{{"Bank code", value[4:8]}, {"Account number", value[8:]}}
	case "BE":
		return []Field{{"Account number", value[4:14]}, {"National check digits", value[14:]}}
	default:
		return nil
	}
}

var ibanCountries = map[string]string{"AT": "Austria", "BE": "Belgium", "DE": "Germany", "ES": "Spain", "FR": "France", "GB": "United Kingdom", "IE": "Ireland", "IT": "Italy", "LU": "Luxembourg", "NL": "Netherlands", "PL": "Poland", "PT": "Portugal", "SE": "Sweden", "CH": "Switzerland"}
var ibanLengths = map[string]int{"AT": 20, "BE": 16, "CH": 21, "DE": 22, "ES": 24, "FR": 27, "GB": 22, "IE": 22, "IT": 27, "LU": 20, "NL": 18, "PL": 28, "PT": 25, "SE": 24}
