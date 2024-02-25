package mblur

import (
	"image"
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
			var pixel FColor
			for j := 0; j < width; j++ {
				pix := FColorFromColor(img.At(offset[j].X+x, offset[j].Y+y))
				for i := 0; i < 3; i++ {
					// multiply each pixel value by their alpha value
					pixel[i] += kernel[j] * pix[i] * pix[3]
				}
				// alpha channel: only multiply by kernel
				pixel[3] += kernel[j] * pix[3]
			}
			if pixel[3] < 1 {
				// we have had some alpha, multiply the various values by gamma
				gamma := PerceptibleReciprocal(pixel[3])
				for i := 0; i < 3; i++ {
					pixel[i] *= gamma
				}
			}
			// this will clamp values
			result.Set(x, y, pixel)
		}
	}
	return result
}
