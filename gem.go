package mblur

import "math"

// https://github.com/ImageMagick/ImageMagick/blob/main/MagickCore/gem.c

// GetOptimalKernelWidth1D computes the optimal kernel radius for a convolution
// filter.  Start with the minimum value of 3 pixels and walk out until we drop
// below the threshold of one pixel numerical accuracy.
func GetOptimalKernelWidth1D(radius, sigma float64) int {
	if radius > MagickEpsilon {
		return int(2.0*math.Ceil(radius) + 1.0)
	}
	gamma := math.Abs(sigma)
	if gamma < MagickEpsilon {
		return 3
	}

	alpha := PerceptibleReciprocal(2 * gamma * gamma)
	beta := PerceptibleReciprocal(MagickSQ2PI * gamma)

	width := 5
	for {
		normalize := 0.0
		j := (width - 1) / 2
		for i := -j; i <= j; i += 1 {
			normalize += math.Exp(-float64(i*i)*alpha) * beta
		}
		value := math.Exp(-float64(j*j)*alpha) * beta / normalize
		if value < QuantumScale || value < MagickEpsilon {
			break
		}
		width += 2
	}
	return width - 2
}
