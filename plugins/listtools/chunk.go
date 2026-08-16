// Package listtools provides local transformations for newline-delimited lists.
package listtools

import (
	"fmt"
	"sort"
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

// PopularItem is one exact list item and its number of occurrences.
type PopularItem struct {
	Item  string
	Count int
}

// FindMostPopular counts non-blank newline-delimited items. Ties are sorted
// alphabetically for deterministic output.
func FindMostPopular(input string) ([]PopularItem, error) {
	counts := make(map[string]int)
	for _, line := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		if strings.TrimSpace(line) != "" {
			counts[line]++
		}
	}
	if len(counts) == 0 {
		return nil, fmt.Errorf("enter at least one non-blank list item")
	}
	items := make([]PopularItem, 0, len(counts))
	for item, count := range counts {
		items = append(items, PopularItem{Item: item, Count: count})
	}
	sort.Slice(items, func(left, right int) bool {
		if items[left].Count != items[right].Count {
			return items[left].Count > items[right].Count
		}
		return items[left].Item < items[right].Item
	})
	return items, nil
}

// UniqueLines returns the first occurrence of each exact non-blank list item.
func UniqueLines(input string) ([]string, error) {
	seen := make(map[string]struct{})
	items := make([]string, 0)
	for _, line := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if _, exists := seen[line]; exists {
			continue
		}
		seen[line] = struct{}{}
		items = append(items, line)
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("enter at least one non-blank list item")
	}
	return items, nil
}

// ReverseLines returns non-blank newline-delimited list items in reverse order.
func ReverseLines(input string) ([]string, error) {
	items := make([]string, 0)
	for _, line := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		if strings.TrimSpace(line) != "" {
			items = append(items, line)
		}
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("enter at least one non-blank list item")
	}
	for left, right := 0, len(items)-1; left < right; left, right = left+1, right-1 {
		items[left], items[right] = items[right], items[left]
	}
	return items, nil
}
