package plugins

import "testing"

func validManifest() Manifest {
	return Manifest{ProtocolVersion: 1, ID: "ports", Name: "Ports", Version: "0.1.0", Pages: []Page{{ID: "overview", Name: "Overview", Default: true}}, Actions: []Action{{ID: "list", Name: "List", Method: "read"}}}
}
func TestValidateManifest(t *testing.T) {
	if err := ValidateManifest(validManifest()); err != nil {
		t.Fatal(err)
	}
	m := validManifest()
	m.ID = "Unsafe ID"
	if err := ValidateManifest(m); err == nil {
		t.Fatal("accepted unsafe ID")
	}
	m = validManifest()
	m.ID = "com.sottey.port-inspector"
	m.Pages[0].ID = "network.overview"
	m.Actions[0].ID = "ports.list"
	if err := ValidateManifest(m); err != nil {
		t.Fatalf("rejected reverse-domain IDs: %v", err)
	}
	m = validManifest()
	m.SearchTerms = []string{"network", "listeners"}
	if err := ValidateManifest(m); err != nil {
		t.Fatalf("rejected search terms: %v", err)
	}
	m = validManifest()
	m.SearchTerms = []string{""}
	if err := ValidateManifest(m); err == nil {
		t.Fatal("accepted blank search term")
	}
	m = validManifest()
	m.ID = ".ports"
	if err := ValidateManifest(m); err == nil {
		t.Fatal("accepted an ID beginning with a period")
	}
	m = validManifest()
	m.ProtocolVersion = 2
	if err := ValidateManifest(m); err == nil {
		t.Fatal("accepted unsupported protocol")
	}
	m = validManifest()
	m.Actions = append(m.Actions, m.Actions[0])
	if err := ValidateManifest(m); err == nil {
		t.Fatal("accepted duplicate action")
	}
}
func TestValidateResponse(t *testing.T) {
	if err := ValidateResponse(Response{OK: true, Data: &View{Components: []Component{{Type: "table"}}}}); err != nil {
		t.Fatal(err)
	}
	if err := ValidateResponse(Response{OK: true, Data: &View{Components: []Component{{Type: "script"}}}}); err == nil {
		t.Fatal("accepted arbitrary component")
	}
	if err := ValidateResponse(Response{OK: true, Data: &View{Components: []Component{{Type: "form", Action: "convert", Submit: "Convert", Fields: []Field{{Name: "json", Label: "JSON", Type: "textarea", Required: true}}}}}}); err != nil {
		t.Fatalf("rejected form component: %v", err)
	}
	if err := ValidateResponse(Response{OK: true, Data: &View{Components: []Component{{Type: "download", Filename: "export.csv", MIMEType: "text/csv", Content: "YSxiCg=="}}}}); err != nil {
		t.Fatalf("rejected download component: %v", err)
	}
}
