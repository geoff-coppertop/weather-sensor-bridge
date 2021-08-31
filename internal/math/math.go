package math

import (
	"math"
)

func Round(val float64, precision uint) float64 {
	rounder := math.Pow10(int(precision))

	return math.Round(val*rounder) / rounder
}
