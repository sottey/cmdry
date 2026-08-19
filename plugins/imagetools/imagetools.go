// Package imagetools provides bounded, in-memory image transformations.
package imagetools

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"strconv"
	"strings"

	"github.com/deepteams/webp"
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
	return EncodePNGWithCompression(source, png.DefaultCompression)
}
func EncodePNGWithCompression(source image.Image, compression png.CompressionLevel) ([]byte, error) {
	var output bytes.Buffer
	err := (&png.Encoder{CompressionLevel: compression}).Encode(&output, source)
	return output.Bytes(), err
}
func EncodeJPEG(source image.Image, quality int) ([]byte, error) {
	var output bytes.Buffer
	err := jpeg.Encode(&output, source, &jpeg.Options{Quality: quality})
	return output.Bytes(), err
}
func EncodeWebP(source image.Image, quality float32) ([]byte, error) {
	var out bytes.Buffer
	err := webp.Encode(&out, source, webp.OptionsForPreset(webp.PresetDefault, quality))
	return out.Bytes(), err
}
func Opacity(source image.Image, percent int) (image.Image, error) {
	if percent < 0 || percent > 100 {
		return nil, fmt.Errorf("opacity must be between 0 and 100")
	}
	b := source.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			r, g, bl, a := source.At(b.Min.X+x, b.Min.Y+y).RGBA()
			out.Set(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(bl >> 8), uint8((a >> 8) * uint32(percent) / 100)})
		}
	}
	return out, nil
}
func parseHex(value string) (color.RGBA, error) {
	value = strings.TrimPrefix(strings.TrimSpace(value), "#")
	if len(value) != 6 {
		return color.RGBA{}, fmt.Errorf("use a six-digit hex color")
	}
	encoded, e := strconv.ParseUint(value, 16, 24)
	if e != nil {
		return color.RGBA{}, fmt.Errorf("use a valid hex color")
	}
	return color.RGBA{uint8(encoded >> 16), uint8(encoded >> 8), uint8(encoded), 255}, nil
}
func ReplaceColor(source image.Image, from, to color.RGBA, tolerance int, transparent bool) (image.Image, error) {
	if tolerance < 0 || tolerance > 255 {
		return nil, fmt.Errorf("tolerance must be from 0 through 255")
	}
	b := source.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			r, g, bl, a := source.At(b.Min.X+x, b.Min.Y+y).RGBA()
			c := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(bl >> 8), uint8(a >> 8)}
			if abs(int(c.R)-int(from.R)) <= tolerance && abs(int(c.G)-int(from.G)) <= tolerance && abs(int(c.B)-int(from.B)) <= tolerance {
				if transparent {
					c.A = 0
				} else {
					c = to
				}
			}
			out.Set(x, y, c)
		}
	}
	return out, nil
}
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
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
