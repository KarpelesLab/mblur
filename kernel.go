package mblur

import (
	"image"
	"image/color"
	"math"
)

type Normalized1DKernel []float64

type Kernel1D []float64

func (k Kernel1D) Normalize() Normalized1DKernel {
	n := make(Normalized1DKernel, len(k))
	var normalize float64
	for i, v := range k {
		n[i] = v
		normalize += v
	}
	for i := range n {
		n[i] /= normalize
	}
	return n
}

func (kernel Kernel1D) Apply(img image.Image, angle float64) image.Image {
	return kernel.Normalize().Apply(img, angle)
}

func (kernel Normalized1DKernel) Apply(img image.Image, angle float64) image.Image {
	width := len(kernel)
	// compute directional offset table
	point := PointInfo{float64(width) * math.Sin(DegreesToRadians(angle)), float64(width) * math.Cos(DegreesToRadians(angle))}
	offset := make([]image.Point, width)
	hypotXY := math.Hypot(point.X, point.Y)
	for w := 0; w < width; w += 1 {
		offset[w].X = int(math.Ceil(float64(w)*point.Y/hypotXY - 0.5))
		offset[w].Y = int(math.Ceil(float64(w)*point.X/hypotXY - 0.5))
	}

	result := CloneImage(img)

	// apply kernel to image at the appropriate angle
	rows := img.Bounds().Dy()
	cols := img.Bounds().Dx()
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			var r, g, b, a float64
			for j := 0; j < width; j++ {
				pr, pg, pb, pa := img.At(offset[j].X+x, offset[j].Y+y).RGBA()
				r += float64(pr) * kernel[j]
				g += float64(pg) * kernel[j]
				b += float64(pb) * kernel[j]
				a += float64(pa) * kernel[j]
			}
			if a < 1 {
				// we have had some alpha, multiply the various values by gamma
				gamma := PerceptibleReciprocal(a)
				r *= gamma
				g *= gamma
				b *= gamma
			}
			// this will clamp values
			pix := color.RGBA64{R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)}
			result.Set(x, y, pix)
		}
	}
	return result
}
