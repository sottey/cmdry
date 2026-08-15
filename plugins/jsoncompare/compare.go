// Package jsoncompare compares two JSON documents structurally.
package jsoncompare

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

const maxDifferences = 500

type Difference struct {
	Kind  string
	Path  string
	Left  string
	Right string
}

type Result struct {
	Differences []Difference
	Truncated   bool
}

// Compare parses exactly one JSON value from each input and compares their
// structure. Object key order is intentionally ignored.
func Compare(leftInput, rightInput string) (Result, error) {
	left, err := decode(leftInput)
	if err != nil {
		return Result{}, fmt.Errorf("parse left JSON: %w", err)
	}
	right, err := decode(rightInput)
	if err != nil {
		return Result{}, fmt.Errorf("parse right JSON: %w", err)
	}
	result := Result{}
	compare("$", left, right, &result)
	return result, nil
}

func decode(input string) (any, error) {
	decoder := json.NewDecoder(strings.NewReader(input))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("input contains more than one JSON value")
		}
		return nil, err
	}
	return value, nil
}

func compare(path string, left, right any, result *Result) {
	if result.Truncated {
		return
	}
	leftObject, leftIsObject := left.(map[string]any)
	rightObject, rightIsObject := right.(map[string]any)
	if leftIsObject || rightIsObject {
		if !leftIsObject || !rightIsObject {
			add(result, Difference{Kind: "changed", Path: path, Left: render(left), Right: render(right)})
			return
		}
		keys := make(map[string]struct{}, len(leftObject)+len(rightObject))
		for key := range leftObject {
			keys[key] = struct{}{}
		}
		for key := range rightObject {
			keys[key] = struct{}{}
		}
		ordered := make([]string, 0, len(keys))
		for key := range keys {
			ordered = append(ordered, key)
		}
		sort.Strings(ordered)
		for _, key := range ordered {
			leftValue, leftOK := leftObject[key]
			rightValue, rightOK := rightObject[key]
			keyPath := objectPath(path, key)
			switch {
			case !leftOK:
				add(result, Difference{Kind: "added", Path: keyPath, Right: render(rightValue)})
			case !rightOK:
				add(result, Difference{Kind: "removed", Path: keyPath, Left: render(leftValue)})
			default:
				compare(keyPath, leftValue, rightValue, result)
			}
		}
		return
	}
	leftArray, leftIsArray := left.([]any)
	rightArray, rightIsArray := right.([]any)
	if leftIsArray || rightIsArray {
		if !leftIsArray || !rightIsArray {
			add(result, Difference{Kind: "changed", Path: path, Left: render(left), Right: render(right)})
			return
		}
		count := len(leftArray)
		if len(rightArray) < count {
			count = len(rightArray)
		}
		for index := 0; index < count; index++ {
			compare(fmt.Sprintf("%s[%d]", path, index), leftArray[index], rightArray[index], result)
		}
		for index := count; index < len(leftArray); index++ {
			add(result, Difference{Kind: "removed", Path: fmt.Sprintf("%s[%d]", path, index), Left: render(leftArray[index])})
		}
		for index := count; index < len(rightArray); index++ {
			add(result, Difference{Kind: "added", Path: fmt.Sprintf("%s[%d]", path, index), Right: render(rightArray[index])})
		}
		return
	}
	if render(left) != render(right) {
		add(result, Difference{Kind: "changed", Path: path, Left: render(left), Right: render(right)})
	}
}

func add(result *Result, difference Difference) {
	if len(result.Differences) >= maxDifferences {
		result.Truncated = true
		return
	}
	result.Differences = append(result.Differences, difference)
}

func objectPath(path, key string) string {
	encoded, _ := json.Marshal(key)
	return path + "[" + string(encoded) + "]"
}

func render(value any) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprint(value)
	}
	return string(encoded)
}
