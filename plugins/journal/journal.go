// Package journal parses the line-delimited JSON output from journalctl.
package journal

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Entry struct {
	Timestamp time.Time
	Priority  int
	Message   string
	Unit      string
}

type journalRecord struct {
	Timestamp string `json:"__REALTIME_TIMESTAMP"`
	Priority  string `json:"PRIORITY"`
	Message   string `json:"MESSAGE"`
	Unit      string `json:"_SYSTEMD_UNIT"`
}

// Parse decodes the JSON Lines output of journalctl --output=json.
func Parse(input []byte) ([]Entry, error) {
	entries := make([]Entry, 0)
	scanner := bufio.NewScanner(bytes.NewReader(input))
	// Messages can be large; retain a bounded but practical maximum line size.
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := bytes.TrimSpace(scanner.Bytes())
		if len(line) == 0 {
			continue
		}
		var record journalRecord
		if err := json.Unmarshal(line, &record); err != nil {
			return nil, fmt.Errorf("decode journal entry %d: %w", lineNumber, err)
		}
		priority, err := strconv.Atoi(record.Priority)
		if err != nil {
			return nil, fmt.Errorf("decode journal entry %d priority: %w", lineNumber, err)
		}
		entry := Entry{Priority: priority, Message: record.Message, Unit: record.Unit}
		if record.Timestamp != "" {
			microseconds, err := strconv.ParseInt(record.Timestamp, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("decode journal entry %d timestamp: %w", lineNumber, err)
			}
			entry.Timestamp = time.UnixMicro(microseconds).Local()
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read journal output: %w", err)
	}
	return entries, nil
}

func (e Entry) Row() map[string]any {
	row := map[string]any{
		"time":     "",
		"priority": PriorityName(e.Priority),
		"unit":     nil,
		"message":  e.Message,
	}
	if !e.Timestamp.IsZero() {
		row["time"] = e.Timestamp.Format("2006-01-02 15:04:05 MST")
	}
	if e.Unit != "" {
		row["unit"] = e.Unit
	}
	return row
}

func PriorityName(priority int) string {
	switch priority {
	case 0:
		return "emergency"
	case 1:
		return "alert"
	case 2:
		return "critical"
	case 3:
		return "error"
	case 4:
		return "warning"
	case 5:
		return "notice"
	case 6:
		return "info"
	case 7:
		return "debug"
	default:
		return "unknown"
	}
}

func Summary(entries []Entry) (errors, warnings int) {
	for _, entry := range entries {
		if entry.Priority <= 3 {
			errors++
		} else if entry.Priority == 4 {
			warnings++
		}
	}
	return errors, warnings
}

// ShortMessage produces a single-line table value while preserving the text.
func ShortMessage(message string) string {
	return strings.ReplaceAll(strings.ReplaceAll(message, "\r", " "), "\n", " ")
}
