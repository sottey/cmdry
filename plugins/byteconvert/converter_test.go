package byteconvert

import "testing"

func TestConvertBinaryAndDecimal(t *testing.T) {
	values, err := Convert(1, "gib")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := values[0], float64(1<<30); got != want {
		t.Fatalf("bytes = %v, want %v", got, want)
	}
	if got, want := values[2], float64(1024); got != want {
		t.Fatalf("MiB = %v, want %v", got, want)
	}
	if got, want := values[8], 1.073741824; got != want {
		t.Fatalf("GB = %v, want %v", got, want)
	}
}

func TestConvertRejectsInvalidInput(t *testing.T) {
	if _, err := Convert(-1, "mb"); err == nil {
		t.Fatal("accepted negative value")
	}
	if _, err := Convert(1, "unknown"); err == nil {
		t.Fatal("accepted unknown unit")
	}
}
