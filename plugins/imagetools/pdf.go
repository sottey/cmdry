package imagetools

import (
	"fmt"
	"math"

	"github.com/raceresult/gopdf"
)

const maxPDFImages = 4
const maxPDFBytes = 6 * 1024 * 1024

// ImagesToPDF creates one A4 page per bounded image and keeps the result in memory.
func ImagesToPDF(contents [][]byte, landscape bool) ([]byte, error) {
	if len(contents) < 1 || len(contents) > maxPDFImages {
		return nil, fmt.Errorf("select from 1 through %d images", maxPDFImages)
	}
	pageWidth, pageHeight := 210.0, 297.0
	if landscape {
		pageWidth, pageHeight = pageHeight, pageWidth
	}
	const margin = 12.0
	builder := gopdf.New()
	for _, source := range contents {
		decoded, _, err := Decode(source)
		if err != nil {
			return nil, err
		}
		img, err := builder.NewImage(source)
		if err != nil {
			return nil, fmt.Errorf("add image to PDF: %w", err)
		}
		bounds := decoded.Bounds()
		scale := math.Min((pageWidth-2*margin)/float64(bounds.Dx()), (pageHeight-2*margin)/float64(bounds.Dy()))
		width, height := float64(bounds.Dx())*scale, float64(bounds.Dy())*scale
		page := builder.NewPage(gopdf.PageSize{gopdf.MM(pageWidth), gopdf.MM(pageHeight)})
		page.AddElement(&gopdf.ImageElement{Img: img, Left: gopdf.MM((pageWidth - width) / 2), Top: gopdf.MM((pageHeight - height) / 2), Width: gopdf.MM(width), Height: gopdf.MM(height)})
	}
	pdf, err := builder.Build()
	if err != nil {
		return nil, fmt.Errorf("build PDF: %w", err)
	}
	if len(pdf) > maxPDFBytes {
		return nil, fmt.Errorf("PDF is too large to download; use fewer or smaller images")
	}
	return pdf, nil
}
