package imagetools

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
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
func RunOpacity() {
	run("image-opacity", "Change Image Opacity", "Adjust one uploaded image opacity locally.", []cmdry.Field{{Name: "image", Label: "Image", Type: "file", Accept: imageAccept, Required: true}, {Name: "opacity", Label: "Opacity (%)", Type: "number", Value: "50", Min: "0", Max: "100", Required: true}}, func(r cmdry.Request) (cmdry.View, error) {
		_, b, e := uploaded(r)
		if e != nil {
			return cmdry.View{}, e
		}
		i, _, e := Decode(b)
		if e != nil {
			return cmdry.View{}, e
		}
		p, e := number(r, "opacity")
		if e != nil {
			return cmdry.View{}, e
		}
		o, e := Opacity(i, p)
		if e != nil {
			return cmdry.View{}, e
		}
		return download("Image opacity", o)
	})
}
func RunWatermark() {
	positions := []cmdry.Option{{Value: "top-left", Label: "Top left"}, {Value: "top-center", Label: "Top center"}, {Value: "top-right", Label: "Top right"}, {Value: "middle-left", Label: "Middle left"}, {Value: "middle-center", Label: "Middle center"}, {Value: "middle-right", Label: "Middle right"}, {Value: "bottom-left", Label: "Bottom left"}, {Value: "bottom-center", Label: "Bottom center"}, {Value: "bottom-right", Label: "Bottom right"}}
	run("watermark-images", "Watermark Images", "Add text to a chosen position on one uploaded image locally.", []cmdry.Field{{Name: "image", Label: "Image", Type: "file", Accept: imageAccept, Required: true}, {Name: "text", Label: "Watermark text", Type: "text", Placeholder: "© Your name", Required: true}, {Name: "position", Label: "Position", Type: "select", Value: "middle-center", Options: positions}, {Name: "color", Label: "Text color", Type: "text", Value: "#ffffff", Required: true}, {Name: "opacity", Label: "Opacity (%)", Type: "number", Value: "50", Min: "1", Max: "100", Required: true}}, func(r cmdry.Request) (cmdry.View, error) {
		_, contents, err := uploaded(r)
		if err != nil {
			return cmdry.View{}, err
		}
		source, _, err := Decode(contents)
		if err != nil {
			return cmdry.View{}, err
		}
		color, err := parseHex(fmt.Sprint(r.Params["color"]))
		if err != nil {
			return cmdry.View{}, err
		}
		opacity, err := number(r, "opacity")
		if err != nil {
			return cmdry.View{}, err
		}
		output, err := Watermark(source, fmt.Sprint(r.Params["text"]), color, opacity, fmt.Sprint(r.Params["position"]))
		if err != nil {
			return cmdry.View{}, err
		}
		return download("Watermarked image", output)
	})
}
func RunImagesToPDF() {
	cmdry.Run(cmdry.Plugin{Manifest: cmdry.Manifest{ProtocolVersion: 1, ID: "com.sottey.images-to-pdf", Name: "Images to PDF", Version: "0.1.0", Description: "Convert up to four uploaded images to a local PDF.", Category: "pdf", Icon: "file", SearchTerms: []string{"image", "pdf", "photos", "pages", "convert"}, Pages: []cmdry.Page{{ID: "overview", Name: "Tool", Default: true, Action: "overview"}}, Permissions: []string{"data.transform"}, Actions: []cmdry.Action{{ID: "overview", Name: "New PDF", Method: "read"}, {ID: "convert", Name: "Create PDF", Method: "write"}}}, Actions: map[string]cmdry.Handler{"overview": func(cmdry.Request) (cmdry.View, error) {
		return cmdry.View{Title: "Images to PDF", Components: []cmdry.Component{{Type: "form", Title: "Create a PDF", Description: "Select up to four PNG, JPEG, or GIF images (4 MiB each). They are held only for this request.", Action: "convert", Submit: "Create PDF", Fields: []cmdry.Field{{Name: "images", Label: "Images", Type: "file", Accept: imageAccept, Multiple: true, Required: true}, {Name: "orientation", Label: "Page orientation", Type: "select", Value: "portrait", Options: []cmdry.Option{{Value: "portrait", Label: "Portrait"}, {Value: "landscape", Label: "Landscape"}}}}}}}, nil
	}, "convert": func(r cmdry.Request) (cmdry.View, error) {
		uploads, contents, err := r.Files("images")
		if err != nil {
			return cmdry.View{}, err
		}
		if len(uploads) > maxPDFImages {
			return cmdry.View{}, fmt.Errorf("select from 1 through %d images", maxPDFImages)
		}
		for _, upload := range uploads {
			if upload.MIMEType != "image/png" && upload.MIMEType != "image/jpeg" && upload.MIMEType != "image/gif" {
				return cmdry.View{}, fmt.Errorf("upload PNG, JPEG, or GIF images")
			}
		}
		landscape := fmt.Sprint(r.Params["orientation"]) == "landscape"
		if value := fmt.Sprint(r.Params["orientation"]); value != "portrait" && value != "landscape" {
			return cmdry.View{}, fmt.Errorf("choose a page orientation")
		}
		pdf, err := ImagesToPDF(contents, landscape)
		if err != nil {
			return cmdry.View{}, err
		}
		return cmdry.View{Title: "Images to PDF", Components: []cmdry.Component{{Type: "metric", Label: "Pages", Value: fmt.Sprint(len(contents))}, {Type: "download", Filename: "cmdry-images.pdf", MIMEType: "application/pdf", Content: base64.StdEncoding.EncodeToString(pdf)}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Create another PDF"}}}}}, nil
	}}})
}
func RunPNGCompression() {
	run("compress-png", "Compress PNG", "Re-encode one uploaded PNG with a selected lossless compression level locally.", []cmdry.Field{{Name: "image", Label: "PNG image", Type: "file", Accept: "image/png", Required: true}, {Name: "compression", Label: "Compression", Type: "select", Value: "default", Options: []cmdry.Option{{Value: "fast", Label: "Fast"}, {Value: "default", Label: "Balanced"}, {Value: "best", Label: "Smallest file"}}}}, func(r cmdry.Request) (cmdry.View, error) {
		upload, b, e := uploaded(r)
		if e != nil {
			return cmdry.View{}, e
		}
		if upload.MIMEType != "image/png" {
			return cmdry.View{}, fmt.Errorf("upload a PNG image")
		}
		i, format, e := Decode(b)
		if e != nil {
			return cmdry.View{}, e
		}
		if format != "png" {
			return cmdry.View{}, fmt.Errorf("upload a PNG image")
		}
		level := png.DefaultCompression
		switch fmt.Sprint(r.Params["compression"]) {
		case "fast":
			level = png.BestSpeed
		case "default":
		case "best":
			level = png.BestCompression
		default:
			return cmdry.View{}, fmt.Errorf("choose a compression level")
		}
		out, e := EncodePNGWithCompression(i, level)
		if e != nil {
			return cmdry.View{}, e
		}
		return encodedDownload("Compressed PNG", "cmdry-image.png", "image/png", out, i.Bounds())
	})
}
func RunColor(id, name string, transparent bool) {
	run(id, name, "Transform a selected PNG/JPEG/GIF color locally.", []cmdry.Field{{Name: "image", Label: "Image", Type: "file", Accept: imageAccept, Required: true}, {Name: "from", Label: "Color to match", Type: "text", Value: "#ffffff", Required: true}, {Name: "to", Label: "Replacement color", Type: "text", Value: "#000000"}, {Name: "tolerance", Label: "Tolerance", Type: "number", Value: "0", Min: "0", Max: "255", Required: true}}, func(r cmdry.Request) (cmdry.View, error) {
		_, b, e := uploaded(r)
		if e != nil {
			return cmdry.View{}, e
		}
		i, _, e := Decode(b)
		if e != nil {
			return cmdry.View{}, e
		}
		f, e := parseHex(fmt.Sprint(r.Params["from"]))
		if e != nil {
			return cmdry.View{}, e
		}
		t, e := parseHex(fmt.Sprint(r.Params["to"]))
		if e != nil {
			return cmdry.View{}, e
		}
		n, e := number(r, "tolerance")
		if e != nil {
			return cmdry.View{}, e
		}
		o, e := ReplaceColor(i, f, t, n, transparent)
		if e != nil {
			return cmdry.View{}, e
		}
		return download(name, o)
	})
}
func RunEncode(id, name, format string) {
	description := "Convert one uploaded image locally."
	if id == "compress-images" {
		description = "Re-encode one uploaded image as a JPEG at a selected quality locally."
	}
	run(id, name, description, []cmdry.Field{{Name: "image", Label: "Image", Type: "file", Accept: imageAccept, Required: true}, {Name: "quality", Label: "Quality", Type: "number", Value: "80", Min: "1", Max: "100", Required: true}}, func(r cmdry.Request) (cmdry.View, error) {
		_, b, e := uploaded(r)
		if e != nil {
			return cmdry.View{}, e
		}
		i, _, e := Decode(b)
		if e != nil {
			return cmdry.View{}, e
		}
		q, e := number(r, "quality")
		if e != nil {
			return cmdry.View{}, e
		}
		if q < 1 || q > 100 {
			return cmdry.View{}, fmt.Errorf("quality must be between 1 and 100")
		}
		var out []byte
		var mime, file string
		if format == "jpg" {
			out, e = EncodeJPEG(i, q)
			mime, file = "image/jpeg", "cmdry-image.jpg"
		} else if format == "webp" {
			out, e = EncodeWebP(i, float32(q))
			mime, file = "image/webp", "cmdry-image.webp"
		} else {
			out, e = EncodePNG(i)
			mime, file = "image/png", "cmdry-image.png"
		}
		if e != nil {
			return cmdry.View{}, e
		}
		return encodedDownload(name, file, mime, out, i.Bounds())
	})
}
func encodedDownload(title, filename, mimeType string, contents []byte, bounds image.Rectangle) (cmdry.View, error) {
	return cmdry.View{Title: title, Components: []cmdry.Component{{Type: "metric", Label: "Output size", Value: fmt.Sprintf("%d × %d", bounds.Dx(), bounds.Dy())}, {Type: "download", Filename: filename, MIMEType: mimeType, Content: base64.StdEncoding.EncodeToString(contents)}, {Type: "actions", Actions: []cmdry.Action{{ID: "overview", Name: "Transform another image"}}}}}, nil
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
	return encodedDownload(title, "cmdry-image.png", "image/png", contents, bounds)
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
