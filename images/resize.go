package images

import (
	"image"

	"golang.org/x/image/draw"
)

func (i *Image) ResizeToSquare(size int) *Image {
	newImg := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.CatmullRom.Scale(newImg, newImg.Rect, i, i.Bounds(), draw.Src, nil)
	return &Image{
		Image:  newImg,
		Format: i.Format,
	}
}
