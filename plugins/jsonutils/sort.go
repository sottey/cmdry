package jsonutils

import (
	"encoding/json"
	"sort"
)

// Sort returns indented JSON with object keys alphabetized. If sortArrays is
// true, every array is also sorted using each item's canonical JSON value.
func Sort(input string, sortArrays bool) ([]byte, error) {
	value, _, err := Parse(input)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(sortValue(value, sortArrays), "", "  ")
}

func sortValue(value any, sortArrays bool) any {
	switch typed := value.(type) {
	case map[string]any:
		output := make(map[string]any, len(typed))
		for key, child := range typed {
			output[key] = sortValue(child, sortArrays)
		}
		return output
	case []any:
		output := make([]any, len(typed))
		for index, child := range typed {
			output[index] = sortValue(child, sortArrays)
		}
		if sortArrays {
			sort.SliceStable(output, func(left, right int) bool {
				leftJSON, _ := json.Marshal(output[left])
				rightJSON, _ := json.Marshal(output[right])
				return string(leftJSON) < string(rightJSON)
			})
		}
		return output
	default:
		return value
	}
}
