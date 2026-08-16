package csvtools

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
)

func writeTable(headers []string, rows [][]string) (string, error) {
	var output bytes.Buffer
	writer := csv.NewWriter(&output)
	if err := writer.Write(headers); err != nil {
		return "", err
	}
	writer.WriteAll(rows)
	if err := writer.Error(); err != nil {
		return "", err
	}
	return output.String(), nil
}

// SwapColumns exchanges two named header columns and their values.
func SwapColumns(input, first, second string) (string, error) {
	headers, rows, err := ReadTable(input)
	if err != nil {
		return "", err
	}
	first, second = strings.TrimSpace(first), strings.TrimSpace(second)
	if first == second {
		return "", fmt.Errorf("choose two different columns")
	}
	a, b := -1, -1
	for index, header := range headers {
		if header == first {
			a = index
		}
		if header == second {
			b = index
		}
	}
	if a < 0 || b < 0 {
		return "", fmt.Errorf("both column names must match CSV headers exactly")
	}
	headers[a], headers[b] = headers[b], headers[a]
	for _, row := range rows {
		row[a], row[b] = row[b], row[a]
	}
	return writeTable(headers, rows)
}

// Transpose turns an input table's rows into output columns.
func Transpose(input string) (string, int, int, error) { return RowsToColumns(input) }
