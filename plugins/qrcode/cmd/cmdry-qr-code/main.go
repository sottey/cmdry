package main

import (
	"encoding/base64"
	"fmt"
	"strconv"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/qrcode"
)

func main() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{
		ProtocolVersion: 1, ID: "com.sottey.qr-code", Name: "QR Code Generator", Version: "0.1.0", Description: "Generate a local PNG QR code from pasted text or a URL.", Category: "image", Icon: "qr",
		SearchTerms: []string{"qr", "qrcode", "barcode", "url", "image", "generator"}, Pages: []cmdry.Page{{ID: "overview", Name: "Generate", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"},
		Actions: []cmdry.Action{{ID: "overview", Name: "New QR code", Method: "read"}, {ID: "generate", Name: "Generate QR code", Method: "write"}},
	}, Actions: map[string]cmdry.Handler{"overview": overview, "generate": generate}})
}

func overview(cmdry.Request) (cmdry.View, error) {
	return cmdry.View{Title: "QR Code Generator", Components: []cmdry.Component{{Type: "form", Title: "Generate QR code", Action: "generate", Submit: "Generate QR code", Description: "Encodes text or a URL locally into a PNG download. Higher recovery levels make a denser code that tolerates more damage.", Fields: []cmdry.Field{{Name: "input", Label: "Text or URL", Type: "textarea", Required: true}, {Name: "recovery", Label: "Error recovery", Type: "select", Value: "medium", Options: []cmdry.Option{{Value: "low", Label: "Low (7%)"}, {Value: "medium", Label: "Medium (15%)"}, {Value: "high", Label: "High (30%)"}}}, {Name: "size", Label: "PNG size (pixels)", Type: "number", Value: "256", Min: "64", Max: "1024", Required: true}}}}}, nil
}

func generate(request cmdry.Request) (cmdry.View, error) {
	size, err := strconv.Atoi(fmt.Sprint(request.Params["size"]))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("PNG size must be a whole number")
	}
	png, err := qrcode.Generate(fmt.Sprint(request.Params["input"]), fmt.Sprint(request.Params["recovery"]), size)
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "QR code ready", Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "QR code generated", Message: "Your PNG was generated locally and is ready to download."}, {Type: "download", Filename: "cmdry-qr-code.png", MIMEType: "image/png", Content: base64.StdEncoding.EncodeToString(png)}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Generate another QR code"}}}}}, nil
}
