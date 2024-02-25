package mblur

import (
	"image"
	"image/color"
	"math"
	"runtime"
	"sync"
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
	procCnt := runtime.NumCPU() // number of processes to run
	var wg sync.WaitGroup
	wg.Add(procCnt)
	for p := 0; p < procCnt; p++ {
		go func(p int) {
			defer wg.Done()
			var r, g, b, a float64
			var gamma float64
			var pix color.RGBA64

			for y := p; y < rows; y += procCnt {
				for x := 0; x < cols; x++ {
					r, g, b, a = 0, 0, 0, 0
					for j := 0; j < width; j++ {
						pr, pg, pb, pa := img.At(offset[j].X+x, offset[j].Y+y).RGBA()
						r += float64(pr) * kernel[j]
						g += float64(pg) * kernel[j]
						b += float64(pb) * kernel[j]
						a += float64(pa) * kernel[j]
					}
					if a < 1 {
						// we have had some alpha, multiply the various values by gamma
						gamma = PerceptibleReciprocal(a)
						r *= gamma
						g *= gamma
						b *= gamma
					}
					// this will clamp values
					pix.R, pix.G, pix.B, pix.A = clamp16(r), clamp16(g), clamp16(b), clamp16(a)
					result.Set(x, y, pix)
				}
			}
		}(p)
	}
	wg.Wait()
	return result
}

func clamp16(v float64) uint16 {
	if math.IsNaN(v) || v <= 0 {
		return 0
	}
	if v >= 0xffff {
		return 0xffff
	}
	return uint16(v)
}
