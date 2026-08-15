// Package byteconvert converts values among binary and decimal byte units.
package byteconvert

import (
	"fmt"
	"math"
)

type Unit struct {
	ID       string
	Label    string
	Bytes    float64
	Standard string
}

var Units = []Unit{
	{ID: "b", Label: "Bytes (B)", Bytes: 1, Standard: "Base"},
	{ID: "kib", Label: "Kibibytes (KiB)", Bytes: 1 << 10, Standard: "Binary"},
	{ID: "mib", Label: "Mebibytes (MiB)", Bytes: 1 << 20, Standard: "Binary"},
	{ID: "gib", Label: "Gibibytes (GiB)", Bytes: 1 << 30, Standard: "Binary"},
	{ID: "tib", Label: "Tebibytes (TiB)", Bytes: 1 << 40, Standard: "Binary"},
	{ID: "pib", Label: "Pebibytes (PiB)", Bytes: 1 << 50, Standard: "Binary"},
	{ID: "kb", Label: "Kilobytes (KB)", Bytes: 1e3, Standard: "Decimal"},
	{ID: "mb", Label: "Megabytes (MB)", Bytes: 1e6, Standard: "Decimal"},
	{ID: "gb", Label: "Gigabytes (GB)", Bytes: 1e9, Standard: "Decimal"},
	{ID: "tb", Label: "Terabytes (TB)", Bytes: 1e12, Standard: "Decimal"},
	{ID: "pb", Label: "Petabytes (PB)", Bytes: 1e15, Standard: "Decimal"},
}

// Convert returns the equivalent value in every supported unit.
func Convert(value float64, inputUnit string) ([]float64, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 {
		return nil, fmt.Errorf("value must be a finite number greater than or equal to zero")
	}
	unit, ok := findUnit(inputUnit)
	if !ok {
		return nil, fmt.Errorf("unsupported unit %q", inputUnit)
	}
	bytes := value * unit.Bytes
	if math.IsInf(bytes, 0) {
		return nil, fmt.Errorf("value is too large to convert")
	}
	results := make([]float64, len(Units))
	for index, target := range Units {
		results[index] = bytes / target.Bytes
	}
	return results, nil
}

func findUnit(id string) (Unit, bool) {
	for _, unit := range Units {
		if unit.ID == id {
			return unit, true
		}
	}
	return Unit{}, false
}
