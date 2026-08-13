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
}
