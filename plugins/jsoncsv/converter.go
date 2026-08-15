// Package jsoncsv converts a JSON object or an array of JSON objects to CSV.
package jsoncsv

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

type Result struct {
	Columns []string
	Rows    []map[string]string
	CSV     []byte
}

// Convert accepts one object or an array of objects. Nested object properties
// use dot paths; nested arrays remain compact JSON inside a single CSV cell.
func Convert(input string) (Result, error) {
	decoder := json.NewDecoder(strings.NewReader(input))
	decoder.UseNumber()
	var root any
	if err := decoder.Decode(&root); err != nil {
		return Result{}, fmt.Errorf("parse JSON: %w", err)
	}
	if err := requireEOF(decoder); err != nil {
		return Result{}, err
	}

	objects, err := records(root)
	if err != nil {
		return Result{}, err
	}
	rows := make([]map[string]string, 0, len(objects))
	columnSet := make(map[string]struct{})
	for index, object := range objects {
		row := make(map[string]string)
		flatten("", object, row)
		if len(row) == 0 {
			return Result{}, fmt.Errorf("record %d contains no exportable fields", index+1)
		}
		for key := range row {
			columnSet[key] = struct{}{}
		}
		rows = append(rows, row)
	}
	columns := make([]string, 0, len(columnSet))
	for key := range columnSet {
		columns = append(columns, key)
	}
	sort.Strings(columns)

	var output bytes.Buffer
	writer := csv.NewWriter(&output)
	if err := writer.Write(columns); err != nil {
		return Result{}, err
	}
	for _, row := range rows {
		record := make([]string, len(columns))
		for index, column := range columns {
			record[index] = row[column]
		}
		if err := writer.Write(record); err != nil {
			return Result{}, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return Result{}, err
	}
	return Result{Columns: columns, Rows: rows, CSV: output.Bytes()}, nil
}

func requireEOF(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return fmt.Errorf("JSON input contains more than one value")
		}
		return fmt.Errorf("parse JSON: %w", err)
	}
	return nil
}

func records(root any) ([]map[string]any, error) {
	switch value := root.(type) {
	case map[string]any:
		return []map[string]any{value}, nil
	case []any:
		if len(value) == 0 {
			return nil, fmt.Errorf("JSON array is empty")
		}
		items := make([]map[string]any, 0, len(value))
		for index, item := range value {
			object, ok := item.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("array item %d is not an object", index+1)
			}
			items = append(items, object)
		}
		return items, nil
	default:
		return nil, fmt.Errorf("JSON must be an object or an array of objects")
	}
}

func flatten(prefix string, value any, row map[string]string) {
	if object, ok := value.(map[string]any); ok {
		keys := make([]string, 0, len(object))
		for key := range object {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			name := key
			if prefix != "" {
				name = prefix + "." + key
			}
			flatten(name, object[key], row)
		}
		return
	}
	if value == nil {
		row[prefix] = ""
		return
	}
	if array, ok := value.([]any); ok {
		encoded, _ := json.Marshal(array)
		row[prefix] = string(encoded)
		return
	}
	row[prefix] = fmt.Sprint(value)
}
