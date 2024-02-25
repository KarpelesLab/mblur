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
func MotionBlurImage(img image.Image, radius, sigma, angle float64) image.Image {
	width := GetOptimalKernelWidth1D(radius, sigma)
	kernel := GetMotionBlurKernel(width, sigma)

	return kernel.Apply(img, angle)
}

func GetMotionBlurKernel(width int, sigma float64) Normalized1DKernel {
	// #define MagickSigma  (fabs(sigma) < MagickEpsilon ? MagickEpsilon : sigma)
	if math.Abs(sigma) < MagickEpsilon {
		sigma = MagickEpsilon
	}
	// Generate a 1-D convolution kernel.
	kernel := make(Normalized1DKernel, width)
	normalize := 0.0
	for i := 0; i < width; i += 1 {
		kernel[i] = (math.Exp((-(float64(i * i)) / (2.0 * sigma * sigma))) / (MagickSQ2PI * sigma))
		normalize += kernel[i]
	}
	for i := 0; i < width; i += 1 {
		kernel[i] /= normalize
	}
	return kernel
}
