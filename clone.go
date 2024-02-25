package mblur

import (
	"image"
	"image/draw"
)

// CloneImage simply copies the given image in a new image.NRGBA
func CloneImage(img image.Image) draw.Image {
	rect := img.Bounds()
	switch img.(type) {
	case *image.Alpha:
		return image.NewAlpha(rect)
	case *image.Alpha16:
		return image.NewAlpha16(rect)
	case *image.Gray:
		return image.NewGray(rect)
	case *image.Gray16:
		return image.NewGray16(rect)
	case *image.NRGBA:
		return image.NewNRGBA(rect)
	case *image.NRGBA64:
		return image.NewNRGBA64(rect)
	case *image.RGBA:
		return image.NewRGBA(rect)
	case *image.RGBA64:
		return image.NewRGBA64(rect)
	default:
		return image.NewNRGBA(rect)
	}
}
