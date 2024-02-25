package mblur_test

import (
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/KarpelesLab/mblur"
)

func TestBlur(t *testing.T) {
	// test the blur
	i, err := loadPngImage("img/gopher.png")
	if err != nil {
		t.Errorf("failed to load img: %s", err)
		return
	}

	res := mblur.MotionBlurImage(i, 0, 4, 45)
	savePngImage("img/gopher-45.png", res)
}

func BenchmarkKernelApply(b *testing.B) {
	kernel := mblur.MotionBlurKernel(0, 4)
	i, err := loadPngImage("img/gopher.png")
	if err != nil {
		b.Errorf("failed to load img: %s", err)
		return
	}

	for n := 0; n < b.N; n++ {
		kernel.Apply(i, 45)
	}
}

func loadPngImage(fn string) (image.Image, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return png.Decode(f)
}

func savePngImage(fn string, img image.Image) error {
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
