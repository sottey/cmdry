package csvtools

import "testing"

func TestRowsToColumns(t *testing.T) {
	output, columns, rows, err := RowsToColumns("name,age\nAda,37\nLin,30\n")
	if err != nil || columns != 2 || rows != 2 || output != "name,Ada,Lin\nage,37,30\n" {
		t.Fatalf("output=%q columns=%d rows=%d err=%v", output, columns, rows, err)
	}
}

func TestCSVToXMLAndYAML(t *testing.T) {
	xmlOutput, rows, err := CSVToXML("first name,age\nAda,37\n")
	if err != nil || rows != 1 || !contains(xmlOutput, "<first_name>Ada</first_name>") {
		t.Fatalf("XML=%q rows=%d err=%v", xmlOutput, rows, err)
	}
	yamlOutput, _, err := CSVToYAML("name,age\nAda,37\n")
	if err != nil || !contains(yamlOutput, "name: Ada") {
		t.Fatalf("YAML=%q err=%v", yamlOutput, err)
	}
}

func contains(input, fragment string) bool {
	return len(input) >= len(fragment) && (input == fragment || stringContains(input, fragment))
}
func stringContains(input, fragment string) bool {
	for index := 0; index+len(fragment) <= len(input); index++ {
		if input[index:index+len(fragment)] == fragment {
			return true
		}
	}
	return false
}
