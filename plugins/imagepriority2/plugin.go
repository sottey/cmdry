package imagepriority2

import (
	"encoding/base64"
	"fmt"
	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
	"github.com/sottey/cmdry/plugins/qrcode"
	"strconv"
)

func RunWiFiQRCode() {
	cmdry.Run(plugin("wifi-qr-code", "Wi-Fi QR Code Generator", "Generate a local QR code for Wi-Fi connection details without retaining credentials.", []string{"wifi", "wi-fi", "qr", "network", "password"}, wifiForm, wifiRun))
}
func RunSVGPlaceholder() {
	cmdry.Run(plugin("svg-placeholder", "SVG Placeholder Generator", "Generate a simple local SVG placeholder image.", []string{"svg", "placeholder", "image", "mockup", "dimensions"}, svgForm, svgRun))
}
func plugin(id, name, description string, terms []string, overview, action cmdry.Handler) cmdry.Plugin {
	return cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey." + id, Name: name, Version: "0.1.0", Description: description, Category: "image", Icon: "image", SearchTerms: terms, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New operation", Method: "read"}, {ID: "run", Name: "Run", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": overview, "run": action}}
}
func wifiForm(cmdry.Request) (cmdry.View, error) {
	return form("Wi-Fi QR Code Generator", "Encodes connection details into a PNG locally. The network name and password are processed only for this request and are not retained.", "Generate Wi-Fi QR code", []cmdry.Field{{Name: "ssid", Label: "Network name (SSID)", Type: "text", Value: "Example Wi-Fi", Required: true}, {Name: "security", Label: "Security", Type: "select", Value: "WPA", Options: []cmdry.Option{{Value: "WPA", Label: "WPA / WPA2 / WPA3"}, {Value: "WEP", Label: "WEP"}, {Value: "nopass", Label: "Open network (no password)"}}}, {Name: "password", Label: "Password (leave blank for open network)", Type: "password"}, {Name: "hidden", Label: "Hidden network", Type: "checkbox"}, {Name: "size", Label: "PNG size (pixels)", Type: "number", Value: "256", Min: "64", Max: "1024", Required: true}})
}
func wifiRun(r cmdry.Request) (cmdry.View, error) {
	payload, err := WiFiPayload(fmt.Sprint(r.Params["ssid"]), fmt.Sprint(r.Params["password"]), fmt.Sprint(r.Params["security"]), fmt.Sprint(r.Params["hidden"]) == "true")
	if err != nil {
		return cmdry.View{}, err
	}
	size, err := strconv.Atoi(fmt.Sprint(r.Params["size"]))
	if err != nil {
		return cmdry.View{}, fmt.Errorf("PNG size must be a whole number")
	}
	png, err := qrcode.Generate(payload, "medium", size)
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "Wi-Fi QR code ready", Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "Complete", Message: "The QR code was generated locally. Connection details are not retained."}, {Type: "download", Filename: "cmdry-wifi-qr-code.png", MIMEType: "image/png", Content: base64.StdEncoding.EncodeToString(png)}, restart("Generate another Wi-Fi QR code")}}, nil
}
func svgForm(cmdry.Request) (cmdry.View, error) {
	return form("SVG Placeholder Generator", "Creates a lightweight SVG locally. Use 3- or 6-digit hex colors.", "Generate SVG", []cmdry.Field{{Name: "width", Label: "Width", Type: "number", Value: "640", Min: "1", Max: "4096", Required: true}, {Name: "height", Label: "Height", Type: "number", Value: "360", Min: "1", Max: "4096", Required: true}, {Name: "label", Label: "Label", Type: "text", Value: "Cmdry placeholder"}, {Name: "background", Label: "Background color", Type: "text", Value: "#1b3543", Required: true}, {Name: "foreground", Label: "Text color", Type: "text", Value: "#ffffff", Required: true}})
}
func svgRun(r cmdry.Request) (cmdry.View, error) {
	parse := func(key string) (int, error) { return strconv.Atoi(fmt.Sprint(r.Params[key])) }
	width, err := parse("width")
	if err != nil {
		return cmdry.View{}, fmt.Errorf("width must be a whole number")
	}
	height, err := parse("height")
	if err != nil {
		return cmdry.View{}, fmt.Errorf("height must be a whole number")
	}
	svg, err := SVGPlaceholder(width, height, fmt.Sprint(r.Params["label"]), fmt.Sprint(r.Params["background"]), fmt.Sprint(r.Params["foreground"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return cmdry.View{Title: "SVG placeholder ready", Components: []cmdry.Component{{Type: "alert", Level: "success", Title: "Complete", Message: "The SVG was generated locally."}, {Type: "download", Filename: "cmdry-placeholder.svg", MIMEType: "image/svg+xml;charset=utf-8", Content: base64.StdEncoding.EncodeToString(svg)}, {Type: "code", Title: "SVG source", Text: string(svg)}, restart("Generate another SVG")}}, nil
}
func form(title, description, submit string, fields []cmdry.Field) (cmdry.View, error) {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "form", Title: title, Description: description + " Your data never leaves this device.", Action: "run", Submit: submit, Fields: fields}}}, nil
}
func restart(name string) cmdry.Component {
	return cmdry.Component{Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: name}}}
}
