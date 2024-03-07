package util

import (
	"math"
)

func TruncateFloat32(value float32, decimalPlaces int) float32 {
	// Use math.Pow to calculate the power of 10
	pow10 := math.Pow10(decimalPlaces)

	// Multiply the value by the power of 10 and round to an integer
	intValue := int(value * float32(pow10))

	// Divide by the power of 10 to get the truncated value
	truncatedValue := float32(intValue) / float32(pow10)

	return truncatedValue
}
