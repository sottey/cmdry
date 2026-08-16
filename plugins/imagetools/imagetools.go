// Package imagetools provides bounded, in-memory image transformations.
package imagetools

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
)

const MaxDimension = 4096

func Decode(contents []byte) (image.Image, string, error) {
	image, format, err := image.Decode(bytes.NewReader(contents))
	if err != nil {
		return nil, "", fmt.Errorf("decode image: %w", err)
	}
	bounds := image.Bounds()
	if bounds.Dx() > MaxDimension || bounds.Dy() > MaxDimension {
		return nil, "", fmt.Errorf("image dimensions must not exceed %d pixels", MaxDimension)
	}
	return image, format, nil
}
func EncodePNG(source image.Image) ([]byte, error) {
	var output bytes.Buffer
	err := png.Encode(&output, source)
	return output.Bytes(), err
}
func EncodeJPEG(source image.Image, quality int) ([]byte, error) {
	var output bytes.Buffer
	err := jpeg.Encode(&output, source, &jpeg.Options{Quality: quality})
	return output.Bytes(), err
}

func Resize(source image.Image, width, height int) (image.Image, error) {
	if width < 1 || height < 1 || width > MaxDimension || height > MaxDimension {
		return nil, fmt.Errorf("width and height must be between 1 and %d", MaxDimension)
	}
	target := image.NewRGBA(image.Rect(0, 0, width, height))
	sourceBounds := source.Bounds()
	for y := 0; y < height; y++ {
		sourceY := sourceBounds.Min.Y + y*sourceBounds.Dy()/height
		for x := 0; x < width; x++ {
			sourceX := sourceBounds.Min.X + x*sourceBounds.Dx()/width
			target.Set(x, y, source.At(sourceX, sourceY))
		}
	}
	return target, nil
}
func Crop(source image.Image, x, y, width, height int) (image.Image, error) {
	bounds := source.Bounds()
	if x < 0 || y < 0 || width < 1 || height < 1 || x+width > bounds.Dx() || y+height > bounds.Dy() {
		return nil, fmt.Errorf("crop must stay within the image bounds (%d × %d)", bounds.Dx(), bounds.Dy())
	}
	target := image.NewRGBA(image.Rect(0, 0, width, height))
	for row := 0; row < height; row++ {
		for column := 0; column < width; column++ {
			target.Set(column, row, source.At(bounds.Min.X+x+column, bounds.Min.Y+y+row))
		}
	}
	return target, nil
}
func Rotate(source image.Image, direction string) (image.Image, error) {
	bounds := source.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	target := image.NewRGBA(image.Rect(0, 0, height, width))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if direction == "clockwise" {
				target.Set(height-1-y, x, source.At(bounds.Min.X+x, bounds.Min.Y+y))
			} else if direction == "counterclockwise" {
				target.Set(y, width-1-x, source.At(bounds.Min.X+x, bounds.Min.Y+y))
			} else {
				return nil, fmt.Errorf("choose clockwise or counterclockwise")
			}
		}
	}
	return target, nil
}
