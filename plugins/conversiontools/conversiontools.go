package conversiontools

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/xuri/excelize/v2"
	"github.com/yuin/goldmark"
)

const (
	maxRows    = 10000
	maxColumns = 100
)

func JSONToTOML(input string) ([]byte, error) {
	decoder := json.NewDecoder(strings.NewReader(input))
	decoder.UseNumber()
	var value map[string]any
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("parse JSON object: %w", err)
	}
	if decoder.More() {
		return nil, fmt.Errorf("input must contain one JSON object")
	}
	var output bytes.Buffer
	if err := toml.NewEncoder(&output).Encode(value); err != nil {
		return nil, fmt.Errorf("encode TOML: %w", err)
	}
	return output.Bytes(), nil
}

func TOMLToJSON(input string) ([]byte, error) {
	var value map[string]any
	if _, err := toml.Decode(input, &value); err != nil {
		return nil, fmt.Errorf("parse TOML: %w", err)
	}
	output, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode JSON: %w", err)
	}
	return append(output, '\n'), nil
}

func FormatTOML(input string) ([]byte, error) {
	var value map[string]any
	if _, err := toml.Decode(input, &value); err != nil {
		return nil, fmt.Errorf("parse TOML: %w", err)
	}
	var output bytes.Buffer
	if err := toml.NewEncoder(&output).Encode(value); err != nil {
		return nil, fmt.Errorf("format TOML: %w", err)
	}
	return output.Bytes(), nil
}

func CSVToXLSX(contents []byte, delimiter rune) ([]byte, int, error) {
	reader := csv.NewReader(bytes.NewReader(contents))
	reader.Comma = delimiter
	records, err := reader.ReadAll()
	if err != nil {
		return nil, 0, fmt.Errorf("read delimited text: %w", err)
	}
	if len(records) == 0 {
		return nil, 0, fmt.Errorf("input contains no rows")
	}
	if len(records) > maxRows {
		return nil, 0, fmt.Errorf("input exceeds %d rows", maxRows)
	}
	book := excelize.NewFile()
	sheet := book.GetSheetName(0)
	for rowIndex, record := range records {
		if len(record) > maxColumns {
			return nil, 0, fmt.Errorf("row %d exceeds %d columns", rowIndex+1, maxColumns)
		}
		for columnIndex, value := range record {
			cell, _ := excelize.CoordinatesToCellName(columnIndex+1, rowIndex+1)
			if err := book.SetCellStr(sheet, cell, spreadsheetText(value)); err != nil {
				return nil, 0, fmt.Errorf("write spreadsheet: %w", err)
			}
		}
	}
	defer book.Close()
	var output bytes.Buffer
	if _, err := book.WriteTo(&output); err != nil {
		return nil, 0, fmt.Errorf("write XLSX: %w", err)
	}
	return output.Bytes(), len(records), nil
}

func XLSXToDelimited(contents []byte, delimiter rune) ([]byte, int, error) {
	book, err := excelize.OpenReader(bytes.NewReader(contents))
	if err != nil {
		return nil, 0, fmt.Errorf("open XLSX: %w", err)
	}
	defer book.Close()
	sheets := book.GetSheetList()
	if len(sheets) == 0 {
		return nil, 0, fmt.Errorf("XLSX contains no worksheets")
	}
	rows, err := book.GetRows(sheets[0])
	if err != nil {
		return nil, 0, fmt.Errorf("read first worksheet: %w", err)
	}
	if len(rows) > maxRows {
		return nil, 0, fmt.Errorf("worksheet exceeds %d rows", maxRows)
	}
	var output bytes.Buffer
	writer := csv.NewWriter(&output)
	writer.Comma = delimiter
	for rowIndex, row := range rows {
		if len(row) > maxColumns {
			return nil, 0, fmt.Errorf("row %d exceeds %d columns", rowIndex+1, maxColumns)
		}
		for index := range row {
			row[index] = spreadsheetText(row[index])
		}
		if err := writer.Write(row); err != nil {
			return nil, 0, fmt.Errorf("write delimited output: %w", err)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, 0, fmt.Errorf("write delimited output: %w", err)
	}
	return output.Bytes(), len(rows), nil
}

func spreadsheetText(value string) string {
	if value != "" && strings.ContainsRune("=+-@", rune(value[0])) {
		return "'" + value
	}
	return value
}

func MarkdownToHTML(input string) ([]byte, error) {
	var output bytes.Buffer
	if err := goldmark.Convert([]byte(input), &output); err != nil {
		return nil, fmt.Errorf("render Markdown: %w", err)
	}
	return output.Bytes(), nil
}

var typeName = regexp.MustCompile(`^[A-Za-z_$][A-Za-z0-9_$]*$`)

func JSONToTypeScript(input, rootName string) (string, error) {
	decoder := json.NewDecoder(strings.NewReader(input))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return "", fmt.Errorf("parse JSON: %w", err)
	}
	if strings.TrimSpace(rootName) == "" {
		rootName = "Root"
	}
	if !typeName.MatchString(rootName) {
		return "", fmt.Errorf("root type name must be a valid TypeScript identifier")
	}
	return "export type " + rootName + " = " + tsType(value, 0) + ";\n", nil
}

func tsType(value any, depth int) string {
	if depth > 12 {
		return "unknown"
	}
	switch typed := value.(type) {
	case nil:
		return "null"
	case bool:
		return "boolean"
	case json.Number:
		return "number"
	case string:
		return "string"
	case []any:
		if len(typed) == 0 {
			return "unknown[]"
		}
		parts := map[string]bool{}
		for _, item := range typed {
			parts[tsType(item, depth+1)] = true
		}
		values := make([]string, 0, len(parts))
		for part := range parts {
			values = append(values, part)
		}
		sort.Strings(values)
		return "(" + strings.Join(values, " | ") + ")[]"
	case map[string]any:
		if len(typed) == 0 {
			return "Record<string, unknown>"
		}
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		parts := make([]string, 0, len(keys))
		for _, key := range keys {
			name := key
			if !typeName.MatchString(name) {
				encoded, _ := json.Marshal(key)
				name = string(encoded)
			}
			parts = append(parts, name+": "+tsType(typed[key], depth+1)+";")
		}
		return "{ " + strings.Join(parts, " ") + " }"
	default:
		return "unknown"
	}
}

func readAll(reader io.Reader) ([]byte, error) { return io.ReadAll(reader) }
