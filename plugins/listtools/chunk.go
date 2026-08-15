// Package listtools provides local transformations for newline-delimited lists.
package listtools

import (
	"fmt"
	"strings"
)

// ChunkLines groups non-blank newline-delimited items into chunks of size.
func ChunkLines(input string, size int) ([][]string, error) {
	if size < 1 {
		return nil, fmt.Errorf("chunk size must be at least 1")
	}
	if size > 10000 {
		return nil, fmt.Errorf("chunk size must not exceed 10,000")
	}
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	items := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			items = append(items, line)
		}
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("enter at least one non-blank list item")
	}
	chunks := make([][]string, 0, (len(items)+size-1)/size)
	for start := 0; start < len(items); start += size {
		end := start + size
		if end > len(items) {
			end = len(items)
		}
		chunks = append(chunks, items[start:end])
	}
	return chunks, nil
}

// FormatChunks joins chunks with an empty line between each group.
func FormatChunks(chunks [][]string) string {
	groups := make([]string, 0, len(chunks))
	for _, chunk := range chunks {
		groups = append(groups, strings.Join(chunk, "\n"))
	}
	return strings.Join(groups, "\n\n")
}
