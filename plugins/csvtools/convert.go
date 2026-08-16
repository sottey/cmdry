package csvtools

import (
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"unicode"

	"go.yaml.in/yaml/v3"
)

// ReadTable parses header-based CSV data and requires complete records.
func ReadTable(input string) ([]string, [][]string, error) {
	reader := csv.NewReader(strings.NewReader(input))
	headers, err := reader.Read()
	if err == io.EOF {
		return nil, nil, fmt.Errorf("CSV input is empty")
	}
	if err != nil {
		return nil, nil, fmt.Errorf("read header: %w", err)
	}
	if len(headers) == 0 {
		return nil, nil, fmt.Errorf("CSV header has no columns")
	}
	seen := make(map[string]struct{}, len(headers))
	for index, header := range headers {
		headers[index] = strings.TrimSpace(header)
		if headers[index] == "" {
			return nil, nil, fmt.Errorf("header column %d is blank", index+1)
		}
		if _, exists := seen[headers[index]]; exists {
			return nil, nil, fmt.Errorf("header %q appears more than once", headers[index])
		}
		seen[headers[index]] = struct{}{}
	}
	rows := make([][]string, 0)
	for rowNumber := 2; ; rowNumber++ {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, fmt.Errorf("read row %d: %w", rowNumber, err)
		}
		if len(row) != len(headers) {
			return nil, nil, fmt.Errorf("row %d has %d fields; expected %d", rowNumber, len(row), len(headers))
		}
		rows = append(rows, row)
	}
	return headers, rows, nil
}

// RowsToColumns transposes a CSV table while retaining its header as the first
// value in each resulting column row.
func RowsToColumns(input string) (string, int, int, error) {
	headers, rows, err := ReadTable(input)
	if err != nil {
		return "", 0, 0, err
	}
	var output bytes.Buffer
	writer := csv.NewWriter(&output)
	for column := range headers {
		record := make([]string, 1, len(rows)+1)
		record[0] = headers[column]
		for _, row := range rows {
			record = append(record, row[column])
		}
		if err := writer.Write(record); err != nil {
			return "", 0, 0, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", 0, 0, err
	}
	return output.String(), len(headers), len(rows), nil
}

func InsertColumns(input string, position int, names []string) (string, error) {
	headers, rows, err := ReadTable(input)
	if err != nil {
		return "", err
	}
	if position < 1 || position > len(headers)+1 || len(names) == 0 {
		return "", fmt.Errorf("choose a valid insertion position and column name")
	}
	for i := range names {
		names[i] = strings.TrimSpace(names[i])
		if names[i] == "" {
			return "", fmt.Errorf("column names cannot be blank")
		}
	}
	at := position - 1
	headers = append(headers[:at], append(names, headers[at:]...)...)
	for i, row := range rows {
		blank := make([]string, len(names))
		rows[i] = append(row[:at], append(blank, row[at:]...)...)
	}
	var out bytes.Buffer
	w := csv.NewWriter(&out)
	_ = w.Write(headers)
	w.WriteAll(rows)
	if e := w.Error(); e != nil {
		return "", e
	}
	return out.String(), nil
}

// CSVToXML renders each data row as a record element with sanitized header
// names used as XML elements.
func CSVToXML(input string) (string, int, error) {
	headers, rows, err := ReadTable(input)
	if err != nil {
		return "", 0, err
	}
	var output bytes.Buffer
	output.WriteString(xml.Header)
	encoder := xml.NewEncoder(&output)
	encoder.Indent("", "  ")
	if err := encoder.EncodeToken(xml.StartElement{Name: xml.Name{Local: "records"}}); err != nil {
		return "", 0, err
	}
	for _, row := range rows {
		if err := encoder.EncodeToken(xml.StartElement{Name: xml.Name{Local: "record"}}); err != nil {
			return "", 0, err
		}
		for index, value := range row {
			if err := encoder.EncodeElement(value, xml.StartElement{Name: xml.Name{Local: xmlName(headers[index])}}); err != nil {
				return "", 0, err
			}
		}
		if err := encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "record"}}); err != nil {
			return "", 0, err
		}
	}
	if err := encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: "records"}}); err != nil {
		return "", 0, err
	}
	if err := encoder.Flush(); err != nil {
		return "", 0, err
	}
	return output.String(), len(rows), nil
}

// CSVToYAML renders each data row as a mapping while preserving cell values as strings.
func CSVToYAML(input string) (string, int, error) {
	headers, rows, err := ReadTable(input)
	if err != nil {
		return "", 0, err
	}
	output := make([]map[string]string, 0, len(rows))
	for _, row := range rows {
		record := make(map[string]string, len(headers))
		for index, header := range headers {
			record[header] = row[index]
		}
		output = append(output, record)
	}
	encoded, err := yaml.Marshal(output)
	return string(encoded), len(rows), err
}

func xmlName(input string) string {
	var output strings.Builder
	for index, value := range input {
		if unicode.IsLetter(value) || value == '_' || (index > 0 && (unicode.IsDigit(value) || value == '-' || value == '.')) {
			output.WriteRune(value)
		} else {
			output.WriteByte('_')
		}
	}
	if output.Len() == 0 {
		return "field"
	}
	return output.String()
}
