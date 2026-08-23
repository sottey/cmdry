package datapriority2

import (
	"strings"
	"testing"
)

func TestYAMLFamilyAndXML(t *testing.T) {
	jsonValue, err := ConvertYAMLFamily("name: Ada\nitems:\n  - one\n", "yaml-json")
	if err != nil || !strings.Contains(string(jsonValue), `"name": "Ada"`) {
		t.Fatalf("json=%q err=%v", jsonValue, err)
	}
	yamlValue, err := ConvertYAMLFamily(`{"name":"Ada"}`, "json-yaml")
	if err != nil || !strings.Contains(string(yamlValue), "name: Ada") {
		t.Fatalf("yaml=%q err=%v", yamlValue, err)
	}
	xmlValue, err := XMLToJSON(`<person id="7"><name>Ada</name><name>Lin</name></person>`)
	if err != nil || !strings.Contains(string(xmlValue), `"@id": "7"`) || !strings.Contains(string(xmlValue), `"name": [`) {
		t.Fatalf("xml=%q err=%v", xmlValue, err)
	}
	if _, err := XMLToJSON(`<one/><two/>`); err == nil {
		t.Fatal("accepted an XML document with multiple root elements")
	}
}

func TestPhoneAndIBAN(t *testing.T) {
	phone, err := ParsePhone("(415) 555-2671", "nanp")
	if err != nil || phone.E164 != "+14155552671" || phone.National != "(415) 555-2671" {
		t.Fatalf("phone=%#v err=%v", phone, err)
	}
	iban, err := ParseIBAN("DE89 3704 0044 0532 0130 00")
	if err != nil || iban.Country != "Germany" || iban.BBAN != "370400440532013000" || len(iban.Fields) != 2 {
		t.Fatalf("iban=%#v err=%v", iban, err)
	}
	if _, err := ParseIBAN("DE89 3704 0044 0532 0130 01"); err == nil {
		t.Fatal("accepted invalid checksum")
	}
	if _, err := ParseIBAN("GB82 WEST 1234"); err == nil {
		t.Fatal("accepted a truncated known-country IBAN")
	}
}
