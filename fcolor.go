package mblur

import (
	"image/color"
	"math"
)

// FColor is a float color where all values are in the [0,1] range
type FColor [4]float64

// FColorFromColor returns a FColor object
func FColorFromColor(col color.Color) FColor {
	switch obj := col.(type) {
	case FColor:
		return obj
	case color.NRGBA:
		// non-alpha-premultiplied
		a := float64(0xff)
		return FColor{float64(obj.R) / a, float64(obj.G) / a, float64(obj.B) / a, float64(obj.A) / a}
	case color.NRGBA64:
		// non-alpha-premultiplied
		a := float64(0xffff)
		return FColor{float64(obj.R) / a, float64(obj.G) / a, float64(obj.B) / a, float64(obj.A) / a}
	default:
		r, g, b, a := col.RGBA()
		if a <= 0 {
			// fully transparent pixel
			return FColor{}
		}
		// TODO rgb would be pre-multiplied by a, we need to undo that
		return FColor{float64(r) / float64(a), float64(g) / float64(a), float64(b) / float64(a), float64(a) / 0xffff}
	}
}

// RGBA returns the alpha-premultiplied red, green, blue and alpha values
func (c FColor) RGBA() (uint32, uint32, uint32, uint32) {
	a := clampAlpha32(c[3], 0xffff) // ensure a is in the 0~0xffff range
	return clampAlpha32(c[0], a), clampAlpha32(c[1], a), clampAlpha32(c[2], a), a
}

// clampAlpha32 returns value of x*a (multiply x by alpha) so that x is adjusted to be in the [0,a] range
func clampAlpha32(x float64, a uint32) uint32 {
	if math.IsNaN(x) || x <= 0.0 {
		return 0
	}
	if x >= 1 {
		// max value
		return a
	}

	return uint32(float64(a) * x)
}
