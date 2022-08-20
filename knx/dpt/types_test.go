// Copyright 2017 Ole Krüger.
// Licensed under the MIT license which can be found in the LICENSE file.

package dpt

import (
	"math"
)

// Define epsilon constant for floating point checks
const epsilon = 1e-3

func abs(x float32) float32 {
	if x < 0.0 {
		return -x
	}
	return x
}

func get_float_quantization_error(value, resolution float32, mantis int) float32 {
	// Calculate the exponent for the value given the mantis and value resolution
	value_m := value / (resolution * float32(mantis))
	value_exp := math.Ceil(math.Log2(float64(value_m)))

	// Calculate the worst quantization error by assuming the
	// mantis to be off by one
	q := math.Pow(2, value_exp)

	// Scale back the quantization error with the given resolution
	return float32(q) / resolution
}
