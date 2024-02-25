package mblur

import (
	"image"
	"image/draw"
)

// CloneImage simply copies the given image in a new image.NRGBA
func CloneImage(img image.Image) *image.NRGBA {
	rect := img.Bounds()
	res := image.NewNRGBA(rect)
	draw.Draw(res, rect, img, image.ZP, draw.Src)
	return res
}
