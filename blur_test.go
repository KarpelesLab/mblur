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
	i, err := loadPngImage("gopher.png")
	if err != nil {
		t.Errorf("failed to load img: %s", err)
		return
	}

	res, _ := mblur.MotionBlurImage(i, 0, 4, 45)
	savePngImage("gopher-45.png", res)
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
