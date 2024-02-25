package mblur

import "math"

type FColor [4]float64

// RGBA returns the alpha-premultiplied red, green, blue and alpha values
func (c FColor) RGBA() (uint32, uint32, uint32, uint32) {
	a := clampFloat(c[3])
	return clampAlpha16(c[0], a), clampAlpha16(c[1], a), clampAlpha16(c[2], a), clampAlpha16(c[3], 1)
}

func clampFloat(x float64) float64 {
	if math.IsNaN(x) || x <= 0.0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}

// clampAlpha16 returns value of x*a (multiply x by alpha) as a 16 bits value in uint32
func clampAlpha16(x, a float64) uint32 {
	x *= a

	if math.IsNaN(x) || x <= 0.0 {
		return 0
	}
	if x > 1 {
		// max value
		return 0xffff
	}

	return uint32(x * 0xffff)
}
