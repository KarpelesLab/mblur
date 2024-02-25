package mblur

import "math"

func GetMotionBlurKernel(width int, sigma float64) []float64 {
	// #define MagickSigma  (fabs(sigma) < MagickEpsilon ? MagickEpsilon : sigma)
	if math.Abs(sigma) < MagickEpsilon {
		sigma = MagickEpsilon
	}
	// Generate a 1-D convolution kernel.
	kernel := make([]float64, width)
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
