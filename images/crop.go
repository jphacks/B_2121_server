package images

import (
	"image"

	"golang.org/x/xerrors"
)

type resizableImage interface {
	SubImage(image.Rectangle) image.Image
}

func getCroppedRectangle(rect image.Rectangle) image.Rectangle {
	width := rect.Dx()
	height := rect.Dy()

	if width == height {
		return rect
	}

	if width > height {
		rect.Min.X += (width - height) / 2
		rect.Max.X = rect.Min.X + height
		return rect
	} else { // height > width
		rect.Min.Y += (height - width) / 2
		rect.Max.Y = rect.Min.Y + width
		return rect
	}
}

func (i *Image) CropToSquare() (*Image, error) {
	resizable, ok := i.Image.(resizableImage)
	if !ok {
		return nil, xerrors.New("failed to crop image")
	}
	rect := i.Bounds()
	newRect := getCroppedRectangle(rect)
	return &Image{
		Image:  resizable.SubImage(newRect),
		Format: i.Format,
	}, nil
}
