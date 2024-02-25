package mblur

import (
	"image"
	"math"
)

// MotionBlurImage simulates motion blur.  We convolve the image with a
// Gaussian operator of the given radius and standard deviation (sigma).
// For reasonable results, radius should be larger than sigma.  Use a
// radius of 0 and MotionBlurImage() selects a suitable radius for you.
// Angle gives the angle of the blurring motion.
func MotionBlurImage(img image.Image, radius, sigma, angle float64) (image.Image, error) {
	width := GetOptimalKernelWidth1D(radius, sigma)
	kernel := GetMotionBlurKernel(width, sigma)

	// compute directional offset table
	point := PointInfo{float64(width) * math.Sin(DegreesToRadians(angle)), float64(width) * math.Cos(DegreesToRadians(angle))}
	offset := make([]image.Point, width)
	for w := 0; w < width; w += 1 {
		offset[w].X = int(math.Ceil(float64(w)*point.Y/math.Hypot(point.X, point.Y) - 0.5))
		offset[w].Y = int(math.Ceil(float64(w)*point.X/math.Hypot(point.X, point.Y) - 0.5))
	}

	// Motion blur image.
	blurImage := CloneImage(img)

	rows := img.Bounds().Dy()
	cols := img.Bounds().Dx()
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			var pixel FColor
			for j := 0; j < width; j++ {
				pix := FColorFromColor(img.At(offset[j].X+x, offset[j].Y+y))
				for i := 0; i < 4; i++ {
					pixel[i] += kernel[j] * pix[i]
				}
			}
			// clamp
			blurImage.Set(x, y, pixel)
		}
	}

	return blurImage, nil
}
