package imagetools

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

func source() image.Image {
	value := image.NewRGBA(image.Rect(0, 0, 2, 3))
	value.Set(0, 0, color.RGBA{R: 255, A: 255})
	value.Set(1, 2, color.RGBA{B: 255, A: 255})
	return value
}

func TestOpacityAndReplaceColor(t *testing.T) {
	value := image.NewRGBA(image.Rect(0, 0, 1, 1))
	value.Set(0, 0, color.RGBA{R: 200, G: 20, B: 10, A: 255})
	opaque, err := Opacity(value, 25)
	if err != nil {
		t.Fatal(err)
	}
	if got := color.RGBAModel.Convert(opaque.At(0, 0)).(color.RGBA); got.A != 63 {
		t.Fatalf("opacity alpha = %d, want 63", got.A)
	}
	from, err := parseHex("#c8140a")
	if err != nil {
		t.Fatal(err)
	}
	to, err := parseHex("0011ff")
	if err != nil {
		t.Fatal(err)
	}
	replaced, err := ReplaceColor(value, from, to, 0, false)
	if err != nil {
		t.Fatal(err)
	}
	if got := color.RGBAModel.Convert(replaced.At(0, 0)).(color.RGBA); got != to {
		t.Fatalf("replacement = %#v, want %#v", got, to)
	}
	transparent, err := ReplaceColor(value, from, to, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	if got := color.RGBAModel.Convert(transparent.At(0, 0)).(color.RGBA); got.A != 0 {
		t.Fatalf("transparent alpha = %d, want 0", got.A)
	}
	if _, err := parseHex("#00gg00"); err == nil {
		t.Fatal("invalid hex color was accepted")
	}
}

func TestEncoders(t *testing.T) {
	value := source()
	pngBytes, err := EncodePNGWithCompression(value, png.BestCompression)
	if err != nil {
		t.Fatal(err)
	}
	if _, format, err := image.Decode(bytes.NewReader(pngBytes)); err != nil || format != "png" {
		t.Fatalf("png format = %q, err = %v", format, err)
	}
	jpegBytes, err := EncodeJPEG(value, 80)
	if err != nil {
		t.Fatal(err)
	}
	if _, format, err := image.Decode(bytes.NewReader(jpegBytes)); err != nil || format != "jpeg" {
		t.Fatalf("jpeg format = %q, err = %v", format, err)
	}
	webpBytes, err := EncodeWebP(value, 80)
	if err != nil {
		t.Fatal(err)
	}
	if len(webpBytes) < 12 || string(webpBytes[:4]) != "RIFF" || string(webpBytes[8:12]) != "WEBP" {
		t.Fatalf("unexpected webp header: %q", webpBytes[:min(len(webpBytes), 12)])
	}
}

func TestWatermarkAndImagesToPDF(t *testing.T) {
	canvas := image.NewRGBA(image.Rect(0, 0, 240, 120))
	for y := 0; y < 120; y++ {
		for x := 0; x < 240; x++ {
			canvas.Set(x, y, color.RGBA{B: 255, A: 255})
		}
	}
	watermarked, err := Watermark(canvas, "Cmdry", color.RGBA{R: 255, A: 255}, 100)
	if err != nil {
		t.Fatal(err)
	}
	changed := false
	for y := 0; y < 120 && !changed; y++ {
		for x := 0; x < 240; x++ {
			if watermarked.At(x, y) != canvas.At(x, y) {
				changed = true
				break
			}
		}
	}
	if !changed {
		t.Fatal("watermark did not change the image")
	}
	contents, err := EncodePNG(canvas)
	if err != nil {
		t.Fatal(err)
	}
	pdf, err := ImagesToPDF([][]byte{contents, contents}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(pdf) < 8 || string(pdf[:5]) != "%PDF-" {
		t.Fatalf("unexpected PDF header: %q", pdf[:min(len(pdf), 8)])
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func TestResizeCropRotate(t *testing.T) {
	resized, err := Resize(source(), 4, 6)
	if err != nil || resized.Bounds().Dx() != 4 || resized.Bounds().Dy() != 6 {
		t.Fatalf("resize=%v %v", resized.Bounds(), err)
	}
	cropped, err := Crop(source(), 1, 2, 1, 1)
	if err != nil || cropped.At(0, 0) != (color.RGBA{B: 255, A: 255}) {
		t.Fatalf("crop=%v %v", cropped.At(0, 0), err)
	}
	rotated, err := Rotate(source(), "clockwise")
	if err != nil || rotated.Bounds().Dx() != 3 || rotated.Bounds().Dy() != 2 {
		t.Fatalf("rotate=%v %v", rotated.Bounds(), err)
	}
}
