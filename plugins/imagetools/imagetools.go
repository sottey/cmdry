// Package imagetools provides bounded, in-memory image transformations.
package imagetools

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strconv"
	"strings"

	"github.com/deepteams/webp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const MaxDimension = 4096
const maxGIFFrames = 120
const maxGIFBytes = 6 * 1024 * 1024

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

// SpeedGIF adjusts an animated GIF's frame delays by the requested rate.
func SpeedGIF(contents []byte, rate float64) ([]byte, int, error) {
	if rate != 0.5 && rate != 1 && rate != 2 && rate != 4 {
		return nil, 0, fmt.Errorf("choose a GIF speed")
	}
	animation, err := decodeGIF(contents)
	if err != nil {
		return nil, 0, err
	}
	for index, delay := range animation.Delay {
		if delay < 1 {
			delay = 10
		}
		adjusted := int(float64(delay)/rate + 0.5)
		if adjusted < 1 {
			adjusted = 1
		}
		animation.Delay[index] = adjusted
	}
	return encodeGIF(animation)
}

// ImagesToGIF turns up to four static images into an endlessly looping GIF.
func ImagesToGIF(contents [][]byte, milliseconds int) ([]byte, int, error) {
	if len(contents) < 1 || len(contents) > maxMultiImages {
		return nil, 0, fmt.Errorf("select from 1 through %d images", maxMultiImages)
	}
	if milliseconds < 10 || milliseconds > 10000 {
		return nil, 0, fmt.Errorf("frame duration must be between 10 and 10,000 milliseconds")
	}
	images := make([]image.Image, len(contents))
	width, height := 0, 0
	for index, source := range contents {
		decoded, _, err := Decode(source)
		if err != nil {
			return nil, 0, err
		}
		images[index] = decoded
		bounds := decoded.Bounds()
		if bounds.Dx() > width {
			width = bounds.Dx()
		}
		if bounds.Dy() > height {
			height = bounds.Dy()
		}
	}
	animation := &gif.GIF{Image: make([]*image.Paletted, len(images)), Delay: make([]int, len(images)), LoopCount: 0}
	for index, source := range images {
		frame := image.NewPaletted(image.Rect(0, 0, width, height), palette.Plan9)
		draw.FloydSteinberg.Draw(frame, frame.Bounds(), source, source.Bounds().Min)
		animation.Image[index] = frame
		animation.Delay[index] = milliseconds / 10
	}
	return encodeGIF(animation)
}

func decodeGIF(contents []byte) (*gif.GIF, error) {
	animation, err := gif.DecodeAll(bytes.NewReader(contents))
	if err != nil {
		return nil, fmt.Errorf("decode GIF: %w", err)
	}
	if len(animation.Image) < 1 || len(animation.Image) > maxGIFFrames {
		return nil, fmt.Errorf("GIF must contain from 1 through %d frames", maxGIFFrames)
	}
	for _, frame := range animation.Image {
		bounds := frame.Bounds()
		if bounds.Dx() > MaxDimension || bounds.Dy() > MaxDimension {
			return nil, fmt.Errorf("GIF dimensions must not exceed %d pixels", MaxDimension)
		}
	}
	return animation, nil
}

func encodeGIF(animation *gif.GIF) ([]byte, int, error) {
	var output bytes.Buffer
	if err := gif.EncodeAll(&output, animation); err != nil {
		return nil, 0, err
	}
	if output.Len() > maxGIFBytes {
		return nil, 0, fmt.Errorf("GIF is too large to download; use fewer or smaller images")
	}
	return output.Bytes(), len(animation.Image), nil
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

// Watermark draws text at a selected grid position over an image without retaining it.
func Watermark(source image.Image, text string, tint color.RGBA, opacity int, position string) (image.Image, error) {
	text = strings.TrimSpace(text)
	if text == "" || len([]rune(text)) > 128 {
		return nil, fmt.Errorf("watermark text must be between 1 and 128 characters")
	}
	if opacity < 1 || opacity > 100 {
		return nil, fmt.Errorf("opacity must be between 1 and 100")
	}
	bounds := source.Bounds()
	out := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(out, out.Bounds(), source, bounds.Min, draw.Src)
	fontData, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return nil, fmt.Errorf("load watermark font: %w", err)
	}
	size := float64(bounds.Dx()) / 10
	if bounds.Dy() < bounds.Dx() {
		size = float64(bounds.Dy()) / 6
	}
	if size < 16 {
		size = 16
	}
	if size > 96 {
		size = 96
	}
	face, err := opentype.NewFace(fontData, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		return nil, fmt.Errorf("create watermark font: %w", err)
	}
	defer face.Close()
	tint.A = uint8(uint32(tint.A) * uint32(opacity) / 100)
	drawer := &font.Drawer{Dst: out, Src: image.NewUniform(tint), Face: face}
	textBounds, _ := drawer.BoundString(text)
	width := (textBounds.Max.X - textBounds.Min.X).Ceil()
	height := face.Metrics().Height.Ceil()
	const margin = 12
	x, y := 0, 0
	switch position {
	case "top-left", "middle-left", "bottom-left":
		x = margin - textBounds.Min.X.Ceil()
	case "top-center", "middle-center", "bottom-center":
		x = (bounds.Dx()-width)/2 - textBounds.Min.X.Ceil()
	case "top-right", "middle-right", "bottom-right":
		x = bounds.Dx() - width - margin - textBounds.Min.X.Ceil()
	default:
		return nil, fmt.Errorf("choose a watermark position")
	}
	switch position {
	case "top-left", "top-center", "top-right":
		y = margin + face.Metrics().Ascent.Ceil()
	case "middle-left", "middle-center", "middle-right":
		y = (bounds.Dy() + face.Metrics().Ascent.Ceil() - face.Metrics().Descent.Ceil()) / 2
	case "bottom-left", "bottom-center", "bottom-right":
		y = bounds.Dy() - margin - (height - face.Metrics().Ascent.Ceil())
	}
	drawer.Dot = fixed.P(x, y)
	drawer.DrawString(text)
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
