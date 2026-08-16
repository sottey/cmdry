package csvjson

import "testing"

func TestConvert(t *testing.T) {
	result, err := Convert("name,age\nAda,37\nLin,\n")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(result.JSON), "[\n  {\n    \"age\": \"37\",\n    \"name\": \"Ada\"\n  },\n  {\n    \"age\": \"\",\n    \"name\": \"Lin\"\n  }\n]\n"; got != want {
		t.Fatalf("JSON = %q, want %q", got, want)
	}
}

func TestConvertRejectsDuplicateHeaders(t *testing.T) {
	if _, err := Convert("name,name\nAda,Lin\n"); err == nil {
		t.Fatal("accepted duplicate headers")
	}
}

func TestConvertDelimitedTSV(t *testing.T) {
	result, err := ConvertDelimited("name\tage\nAda\t37\n", '\t')
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(result.JSON), "[\n  {\n    \"age\": \"37\",\n    \"name\": \"Ada\"\n  }\n]\n"; got != want {
		t.Fatalf("JSON = %q, want %q", got, want)
	}
}
