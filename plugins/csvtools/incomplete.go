package csvtools

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// IncompleteRecord identifies one CSV data row that has empty fields.
type IncompleteRecord struct {
	Row     int
	Columns []string
}

// FindIncompleteRecords finds empty values in comma-separated CSV data. The
// first row is used as the header, and short rows count as missing trailing
// columns.
func FindIncompleteRecords(input string) ([]IncompleteRecord, int, error) {
	reader := csv.NewReader(strings.NewReader(input))
	reader.FieldsPerRecord = -1
	headers, err := reader.Read()
	if err == io.EOF {
		return nil, 0, fmt.Errorf("CSV input is empty")
	}
	if err != nil {
		return nil, 0, fmt.Errorf("read header: %w", err)
	}
	if len(headers) == 0 {
		return nil, 0, fmt.Errorf("CSV header has no columns")
	}
	for index, header := range headers {
		if strings.TrimSpace(header) == "" {
			headers[index] = fmt.Sprintf("Column %d", index+1)
		}
	}
	records, dataRows := make([]IncompleteRecord, 0), 0
	for rowNumber := 2; ; rowNumber++ {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, 0, fmt.Errorf("read row %d: %w", rowNumber, err)
		}
		dataRows++
		missing := make([]string, 0)
		for index, header := range headers {
			if index >= len(row) || row[index] == "" {
				missing = append(missing, header)
			}
		}
		if len(missing) > 0 {
			records = append(records, IncompleteRecord{Row: rowNumber, Columns: missing})
		}
	}
	return records, dataRows, nil
}
