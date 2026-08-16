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
