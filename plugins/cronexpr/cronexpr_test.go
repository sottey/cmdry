package cronexpr

import "testing"

func TestExplainStandardExpression(t *testing.T) {
	fields, notice, err := Explain("*/15 9-17 * JAN,MAR MON-FRI")
	if err != nil || len(fields) != 5 || fields[0].Explanation != "every 15 minutes" || fields[4].Explanation != "days of the week MON through FRI" || notice == "" {
		t.Fatalf("fields = %#v, notice = %q, err = %v", fields, notice, err)
	}
}

func TestExplainShortcut(t *testing.T) {
	fields, description, err := Explain("@daily")
	if err != nil || fields != nil || description != "At midnight every day." {
		t.Fatalf("fields = %#v, description = %q, err = %v", fields, description, err)
	}
}

func TestExplainRejectsInvalidExpressions(t *testing.T) {
	if _, _, err := Explain("0 25 * * *"); err == nil {
		t.Fatal("expected invalid hour error")
	}
}
