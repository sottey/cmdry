package xmlutils

import "testing"

func TestValidate(t *testing.T) {
	info, err := Validate(`<root language="en"><child/></root>`)
	if err != nil || info.Root != "root" || info.Attributes != 1 {
		t.Fatalf("info = %#v, err = %v", info, err)
	}
}

func TestValidateReturnsLineForInvalidXML(t *testing.T) {
	if _, err := Validate("<root>\n<child></root>"); err == nil {
		t.Fatal("expected invalid XML error")
	}
}
