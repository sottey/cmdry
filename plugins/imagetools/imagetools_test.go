package imagetools

import (
	"image"
	"image/color"
	"testing"
)

func source() image.Image {
	value := image.NewRGBA(image.Rect(0, 0, 2, 3))
	value.Set(0, 0, color.RGBA{R: 255, A: 255})
	value.Set(1, 2, color.RGBA{B: 255, A: 255})
	return value
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
