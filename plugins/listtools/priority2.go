package listtools

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

func priorityItems(input string) ([]string, error) {
	items := make([]string, 0)
	for _, line := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		if strings.TrimSpace(line) != "" {
			items = append(items, line)
		}
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("enter at least one non-blank list item")
	}
	return items, nil
}

// RotateLines moves list items left by positions; negative values rotate right.
func RotateLines(input string, positions int) ([]string, error) {
	items, err := priorityItems(input)
	if err != nil {
		return nil, err
	}
	shift := positions % len(items)
	if shift < 0 {
		shift += len(items)
	}
	return append(items[shift:], items[:shift]...), nil
}

// ShuffleLines returns a cryptographically shuffled copy of non-blank items.
func ShuffleLines(input string) ([]string, error) {
	items, err := priorityItems(input)
	if err != nil {
		return nil, err
	}
	for index := len(items) - 1; index > 0; index-- {
		pick, err := rand.Int(rand.Reader, big.NewInt(int64(index+1)))
		if err != nil {
			return nil, fmt.Errorf("shuffle list: %w", err)
		}
		items[index], items[pick.Int64()] = items[pick.Int64()], items[index]
	}
	return items, nil
}

// WrapLines adds prefix and suffix to each non-blank input line.
func WrapLines(input, prefix, suffix string) ([]string, error) {
	items, err := priorityItems(input)
	if err != nil {
		return nil, err
	}
	for index, item := range items {
		items[index] = prefix + item + suffix
	}
	return items, nil
}

// UnwrapLines removes one matching prefix and suffix from every item.
func UnwrapLines(input, prefix, suffix string) ([]string, error) {
	if prefix == "" && suffix == "" {
		return nil, fmt.Errorf("enter a prefix or suffix to remove")
	}
	items, err := priorityItems(input)
	if err != nil {
		return nil, err
	}
	for index, item := range items {
		if prefix != "" {
			if !strings.HasPrefix(item, prefix) {
				return nil, fmt.Errorf("item %d does not start with the prefix", index+1)
			}
			item = strings.TrimPrefix(item, prefix)
		}
		if suffix != "" {
			if !strings.HasSuffix(item, suffix) {
				return nil, fmt.Errorf("item %d does not end with the suffix", index+1)
			}
			item = strings.TrimSuffix(item, suffix)
		}
		items[index] = item
	}
	return items, nil
}
