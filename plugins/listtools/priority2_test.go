package listtools

import "testing"

func TestPriority2ListTools(t *testing.T) {
	rotated, err := RotateLines("a\nb\nc", -1)
	if err != nil || joined(rotated) != "c\na\nb" {
		t.Fatalf("rotate=%q err=%v", joined(rotated), err)
	}
	wrapped, err := WrapLines("a\na", "[", "]")
	if err != nil || joined(wrapped) != "[a]\n[a]" {
		t.Fatalf("wrap=%q err=%v", joined(wrapped), err)
	}
	unwrapped, err := UnwrapLines("[a]\n[b]", "[", "]")
	if err != nil || joined(unwrapped) != "a\nb" {
		t.Fatalf("unwrap=%q err=%v", joined(unwrapped), err)
	}
	shuffled, err := ShuffleLines("a\nb\nc")
	if err != nil || len(shuffled) != 3 {
		t.Fatalf("shuffle=%q err=%v", joined(shuffled), err)
	}
}

func joined(items []string) string {
	result := ""
	for index, item := range items {
		if index > 0 {
			result += "\n"
		}
		result += item
	}
	return result
}
