package neat

import (
	"math"

	"github.com/wowsims/wotlk/sim/core"
)

func MinFloat64Slice(slice []float64) float64 {
	min := slice[0]
	for _, v := range slice {
		if v <= min {
			min = v
		}
	}
	return min
}

func MaxFloat64Slice(slice []float64) float64 {
	max := slice[0]
	for _, v := range slice {
		if v >= max {
			max = v
		}
	}
	return max
}

func Regularize(slice []float64) []float64 {
	max := MaxFloat64Slice(slice)
	min := MinFloat64Slice(slice)
	factor := core.TernaryFloat64(math.Abs(max) >= math.Abs(min), math.Abs(max), math.Abs(min))
	for i := range slice {
		slice[i] /= factor
	}
	return slice
}
