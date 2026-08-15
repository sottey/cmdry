// Package csvtools provides local CSV transformations.
package csvtools

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

var Separators = map[string]rune{"comma": ',', "semicolon": ';', "tab": '\t', "pipe": '|'}

// ChangeSeparator parses CSV with source and rewrites it with destination.
// Quoted fields and embedded separators are preserved by encoding/csv.
func ChangeSeparator(input, source, destination string) (string, int, error) {
	sourceRune, ok := Separators[source]
	if !ok {
		return "", 0, fmt.Errorf("unsupported source separator")
	}
	destinationRune, ok := Separators[destination]
	if !ok {
		return "", 0, fmt.Errorf("unsupported destination separator")
	}
	reader := csv.NewReader(strings.NewReader(input))
	reader.Comma = sourceRune
	reader.FieldsPerRecord = -1
	rows := make([][]string, 0)
	for rowNumber := 1; ; rowNumber++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", 0, fmt.Errorf("read row %d: %w", rowNumber, err)
		}
		rows = append(rows, record)
	}
	if len(rows) == 0 {
		return "", 0, fmt.Errorf("CSV input is empty")
	}
	var output bytes.Buffer
	writer := csv.NewWriter(&output)
	writer.Comma = destinationRune
	if err := writer.WriteAll(rows); err != nil {
		return "", 0, fmt.Errorf("write CSV: %w", err)
	}
	return output.String(), len(rows), nil
}
