// Package csvjson converts header-based CSV data to JSON records.
package csvjson

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Result struct {
	Headers []string
	Rows    []map[string]string
	JSON    []byte
}

// Convert accepts CSV whose first row supplies unique object property names.
func Convert(input string) (Result, error) {
	reader := csv.NewReader(strings.NewReader(input))
	reader.FieldsPerRecord = -1
	headers, err := reader.Read()
	if err == io.EOF {
		return Result{}, fmt.Errorf("CSV input is empty")
	}
	if err != nil {
		return Result{}, fmt.Errorf("read header: %w", err)
	}
	if len(headers) == 0 {
		return Result{}, fmt.Errorf("CSV header has no columns")
	}
	seen := make(map[string]struct{}, len(headers))
	for index, header := range headers {
		headers[index] = strings.TrimSpace(header)
		if headers[index] == "" {
			return Result{}, fmt.Errorf("header column %d is blank", index+1)
		}
		if _, exists := seen[headers[index]]; exists {
			return Result{}, fmt.Errorf("header %q appears more than once", headers[index])
		}
		seen[headers[index]] = struct{}{}
	}
	rows := make([]map[string]string, 0)
	for rowNumber := 2; ; rowNumber++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return Result{}, fmt.Errorf("read row %d: %w", rowNumber, err)
		}
		if len(record) != len(headers) {
			return Result{}, fmt.Errorf("row %d has %d fields; expected %d", rowNumber, len(record), len(headers))
		}
		row := make(map[string]string, len(headers))
		for index, header := range headers {
			row[header] = record[index]
		}
		rows = append(rows, row)
	}
	encoded, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return Result{}, err
	}
	return Result{Headers: headers, Rows: rows, JSON: append(bytes.Clone(encoded), '\n')}, nil
}
