package numbertools

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// GenerateRandomDecimals returns cryptographically secure fixed-precision
// decimal values within inclusive decimal bounds.
func GenerateRandomDecimals(minInput, maxInput string, places, count int) ([]string, error) {
	if places < 0 || places > 9 {
		return nil, fmt.Errorf("decimal places must be between 0 and 9")
	}
	if count < 1 || count > 10000 {
		return nil, fmt.Errorf("count must be between 1 and 10,000")
	}
	scale := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(places)), nil)
	min, err := scaledDecimal(minInput, scale)
	if err != nil {
		return nil, fmt.Errorf("minimum: %w", err)
	}
	max, err := scaledDecimal(maxInput, scale)
	if err != nil {
		return nil, fmt.Errorf("maximum: %w", err)
	}
	if min.Cmp(max) > 0 {
		return nil, fmt.Errorf("minimum must not exceed maximum")
	}
	width := new(big.Int).Sub(max, min)
	width.Add(width, big.NewInt(1))
	values := make([]string, 0, count)
	for range count {
		offset, err := rand.Int(rand.Reader, width)
		if err != nil {
			return nil, fmt.Errorf("generate random number: %w", err)
		}
		value := new(big.Int).Add(min, offset)
		values = append(values, formatScaledDecimal(value, places))
	}
	return values, nil
}

func scaledDecimal(input string, scale *big.Int) (*big.Int, error) {
	value, ok := new(big.Rat).SetString(strings.TrimSpace(input))
	if !ok {
		return nil, fmt.Errorf("must be a decimal number")
	}
	value.Mul(value, new(big.Rat).SetInt(scale))
	if value.Denom().Cmp(big.NewInt(1)) != 0 {
		return nil, fmt.Errorf("has more decimal places than configured")
	}
	return new(big.Int).Set(value.Num()), nil
}

func formatScaledDecimal(value *big.Int, places int) string {
	if places == 0 {
		return value.String()
	}
	sign, digits := "", value.String()
	if strings.HasPrefix(digits, "-") {
		sign, digits = "-", digits[1:]
	}
	if len(digits) <= places {
		digits = strings.Repeat("0", places-len(digits)+1) + digits
	}
	return sign + digits[:len(digits)-places] + "." + digits[len(digits)-places:]
}
