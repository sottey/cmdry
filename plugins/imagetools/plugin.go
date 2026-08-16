package imagetools

import (
	"encoding/base64"
	"fmt"
	"image"
	"strconv"

	cmdry "github.com/sottey/cmdry/plugin-sdk/go"
)

const imageAccept = "image/png,image/jpeg,image/gif"

func RunResize() {
	run("resize-image", "Resize Image", "Resize one uploaded image locally; it is held only for this request.", []cmdry.Field{{Name: "image", Label: "Image", Type: "file", Accept: imageAccept, Required: true}, {Name: "width", Label: "Width (pixels)", Type: "number", Required: true}, {Name: "height", Label: "Height (pixels)", Type: "number", Required: true}}, resizeAction)
}
func RunCrop() {
	run("crop-image", "Crop Image", "Crop one uploaded image locally; it is held only for this request.", []cmdry.Field{{Name: "image", Label: "Image", Type: "file", Accept: imageAccept, Required: true}, {Name: "x", Label: "Left (pixels)", Type: "number", Value: "0", Min: "0", Required: true}, {Name: "y", Label: "Top (pixels)", Type: "number", Value: "0", Min: "0", Required: true}, {Name: "width", Label: "Width (pixels)", Type: "number", Required: true}, {Name: "height", Label: "Height (pixels)", Type: "number", Required: true}}, cropAction)
}
func RunRotate() {
	run("rotate-image", "Rotate Image", "Rotate one uploaded image 90 degrees locally; it is held only for this request.", []cmdry.Field{{Name: "image", Label: "Image", Type: "file", Accept: imageAccept, Required: true}, {Name: "direction", Label: "Direction", Type: "select", Value: "clockwise", Options: []cmdry.Option{{Value: "clockwise", Label: "Clockwise"}, {Value: "counterclockwise", Label: "Counterclockwise"}}}}, rotateAction)
}

func run(id, name, description string, fields []cmdry.Field, action cmdry.Handler) {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey." + id, Name: name, Version: "0.1.0", Description: description, Category: "image", Icon: "image", SearchTerms: []string{"image", "upload", "local"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New image", Method: "read"}, {ID: "run", Name: "Transform image", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": func(cmdry.Request) (cmdry.View, error) {
		return cmdry.View{Title: name, Components: []cmdry.Component{{Type: "form", Title: name, Description: "Images are limited to 4 MiB, processed in memory, and never written to host paths.", Action: "run", Submit: "Transform image", Fields: fields}}}, nil
	}, "run": action}})
}
func uploaded(request cmdry.Request) (cmdry.Upload, []byte, error) {
	upload, contents, err := request.File("image")
	if err != nil {
		return cmdry.Upload{}, nil, err
	}
	if upload.MIMEType != "image/png" && upload.MIMEType != "image/jpeg" && upload.MIMEType != "image/gif" {
		return cmdry.Upload{}, nil, fmt.Errorf("upload a PNG, JPEG, or GIF image")
	}
	return upload, contents, nil
}
func number(request cmdry.Request, key string) (int, error) {
	value, err := strconv.Atoi(fmt.Sprint(request.Params[key]))
	if err != nil {
		return 0, fmt.Errorf("%s must be a whole number", key)
	}
	return value, nil
}
func download(title string, source image.Image) (cmdry.View, error) {
	contents, err := EncodePNG(source)
	if err != nil {
		return cmdry.View{}, err
	}
	bounds := source.Bounds()
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "metric", Label: "Output size", Value: fmt.Sprintf("%d × %d", bounds.Dx(), bounds.Dy())}, {Type: "download", Filename: "cmdry-image.png", MIMEType: "image/png", Content: base64.StdEncoding.EncodeToString(contents)}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Transform another image"}}}}}, nil
}
func resizeAction(request cmdry.Request) (cmdry.View, error) {
	_, contents, err := uploaded(request)
	if err != nil {
		return cmdry.View{}, err
	}
	source, _, err := Decode(contents)
	if err != nil {
		return cmdry.View{}, err
	}
	width, err := number(request, "width")
	if err != nil {
		return cmdry.View{}, err
	}
	height, err := number(request, "height")
	if err != nil {
		return cmdry.View{}, err
	}
	output, err := Resize(source, width, height)
	if err != nil {
		return cmdry.View{}, err
	}
	return download("Resized image", output)
}
func cropAction(request cmdry.Request) (cmdry.View, error) {
	_, contents, err := uploaded(request)
	if err != nil {
		return cmdry.View{}, err
	}
	source, _, err := Decode(contents)
	if err != nil {
		return cmdry.View{}, err
	}
	x, err := number(request, "x")
	if err != nil {
		return cmdry.View{}, err
	}
	y, err := number(request, "y")
	if err != nil {
		return cmdry.View{}, err
	}
	width, err := number(request, "width")
	if err != nil {
		return cmdry.View{}, err
	}
	height, err := number(request, "height")
	if err != nil {
		return cmdry.View{}, err
	}
	output, err := Crop(source, x, y, width, height)
	if err != nil {
		return cmdry.View{}, err
	}
	return download("Cropped image", output)
}
func rotateAction(request cmdry.Request) (cmdry.View, error) {
	_, contents, err := uploaded(request)
	if err != nil {
		return cmdry.View{}, err
	}
	source, _, err := Decode(contents)
	if err != nil {
		return cmdry.View{}, err
	}
	output, err := Rotate(source, fmt.Sprint(request.Params["direction"]))
	if err != nil {
		return cmdry.View{}, err
	}
	return download("Rotated image", output)
}
