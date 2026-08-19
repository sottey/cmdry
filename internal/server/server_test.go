package server

import (
	"bytes"
	"encoding/base64"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"github.com/sottey/cmdry/internal/plugins"
)

func TestActionParamsReadsEphemeralUpload(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	_ = writer.WriteField("width", "12")
	part, err := writer.CreateFormFile("image", "example.png")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write([]byte("\x89PNG\r\n\x1a\ncontent"))
	_ = writer.Close()
	req := httptest.NewRequest("POST", "/", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	params, err := actionParams(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	if params["width"] != "12" {
		t.Fatalf("width=%#v", params["width"])
	}
	upload, ok := params["image"].(plugins.Upload)
	if !ok {
		t.Fatalf("upload=%T", params["image"])
	}
	contents, err := base64.StdEncoding.DecodeString(upload.Content)
	if err != nil || string(contents) != "\x89PNG\r\n\x1a\ncontent" {
		t.Fatalf("contents=%q err=%v", contents, err)
	}
}

func TestMakeGroupsKeepsServerLast(t *testing.T) {
	entries := []plugins.Registered{
		{Manifest: plugins.Manifest{ID: "com.sottey.ports", Name: "Port Inspector", Category: "server"}},
		{Manifest: plugins.Manifest{ID: "com.sottey.base64", Name: "Base64", Category: "text"}},
		{Manifest: plugins.Manifest{ID: "com.sottey.processes", Name: "Processes", Category: "server"}},
		{Manifest: plugins.Manifest{ID: "com.sottey.sum", Name: "Sum", Category: "number"}},
	}
	groups := makeGroups(entries, nil)
	if len(groups) != 3 {
		t.Fatalf("group count = %d, want 3", len(groups))
	}
	if groups[0].ID != "text" || groups[1].ID != "number" || groups[2].ID != "server" {
		t.Fatalf("group order = %q, %q, %q", groups[0].ID, groups[1].ID, groups[2].ID)
	}
	if len(groups[2].Plugins) != 2 {
		t.Fatalf("server plugin count = %d, want 2", len(groups[2].Plugins))
	}
}

func TestActionParamsReadsMultipleEphemeralUploads(t *testing.T) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for _, name := range []string{"one.png", "two.png"} {
		part, err := writer.CreateFormFile("images", name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := part.Write([]byte("image-" + name)); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("POST", "/", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	params, err := actionParams(httptest.NewRecorder(), req)
	if err != nil {
		t.Fatal(err)
	}
	uploads, ok := params["images"].([]plugins.Upload)
	if !ok || len(uploads) != 2 {
		t.Fatalf("uploads=%#v", params["images"])
	}
	if uploads[0].Name != "one.png" || uploads[1].Name != "two.png" {
		t.Fatalf("upload names=%q, %q", uploads[0].Name, uploads[1].Name)
	}
}
