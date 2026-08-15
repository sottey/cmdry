package cryptotools

import "testing"

func TestHashAndVerify(t *testing.T) {
	hash, err := Hash("correct horse battery staple", 4)
	if err != nil {
		t.Fatal(err)
	}
	if matched, err := Verify("correct horse battery staple", hash); err != nil || !matched {
		t.Fatalf("matched = %v, err = %v", matched, err)
	}
	if matched, err := Verify("incorrect", hash); err != nil || matched {
		t.Fatalf("matched = %v, err = %v", matched, err)
	}
}

func TestHashRejectsInvalidCost(t *testing.T) {
	if _, err := Hash("text", 3); err == nil {
		t.Fatal("expected invalid cost error")
	}
}
