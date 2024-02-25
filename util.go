package mblur

import "math"

const (
	MagickEpsilon = 1.0e-12
	MagickSQ2PI   = 2.50662827463100024161235523934010416269302368164062
	QuantumRange  = float64(18446744073709551615.0)
	QuantumScale  = 1 / QuantumRange
)

// PerceptibleReciprocal returns 1/x where x is perceptible (not unlimited or infinitesimal).
func PerceptibleReciprocal(x float64) float64 {
	sign := 1.0
	if x < 0.0 {
		sign = -1.0
	}
	if (sign * x) > MagickEpsilon {
		return 1 / x
	}
	return sign / MagickEpsilon
}

func DegreesToRadians(deg float64) float64 {
	return math.Pi * deg / 180
}

func ClampToQuantum(quantum float64) float64 {
	if math.IsNaN(quantum) || quantum <= 0.0 {
		return 0
	}
	if quantum >= QuantumRange {
		return QuantumRange
	}
	return quantum + 0.5
}
