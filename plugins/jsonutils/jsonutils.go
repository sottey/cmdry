// Package jsonutils provides local JSON validation and transformation helpers.
package jsonutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Parse accepts exactly one JSON value and returns both its decoded value and
// compact JSON representation.
func Parse(input string) (any, []byte, error) {
	decoder := json.NewDecoder(strings.NewReader(input))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, nil, err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, nil, fmt.Errorf("input contains more than one JSON value")
		}
		return nil, nil, err
	}
	var compact bytes.Buffer
	if err := json.Compact(&compact, []byte(input)); err != nil {
		return nil, nil, err
	}
	return value, compact.Bytes(), nil
}

func RootType(value any) string {
	switch value.(type) {
	case map[string]any:
		return "object"
	case []any:
		return "array"
	case string:
		return "string"
	case json.Number:
		return "number"
	case bool:
		return "boolean"
	case nil:
		return "null"
	default:
		return "value"
	}
}

// Stringify encodes a compact, valid JSON document as a JSON string literal.
func Stringify(input string) ([]byte, error) {
	_, compact, err := Parse(input)
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(compact))
}

// EscapeString encodes arbitrary text as one JSON string literal. Unlike
// Stringify, the input is not required to be valid JSON.
func EscapeString(input string) ([]byte, error) {
	return json.Marshal(input)
}
