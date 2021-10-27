package images

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"golang.org/x/xerrors"
)

type Image struct {
	image.Image
	Format string
}

func LoadImage(reader io.Reader) (*Image, error) {
	img, format, err := image.Decode(reader)
	if err != nil {
		return nil, xerrors.Errorf("failed to load image: %w", err)
	}

	return &Image{
		Image:  img,
		Format: format,
	}, nil
}

func (i *Image) Save(writer io.Writer) error {
	var err error
	switch i.Format {
	case "png":
		err = png.Encode(writer, i.Image)
	case "jpeg":
		err = jpeg.Encode(writer, i.Image, &jpeg.Options{Quality: 75})
	default:
		return xerrors.Errorf("unknown format: %s", i.Format)
	}
	if err != nil {
		return xerrors.Errorf("failed to save image: %w", err)
	}
	return nil
}

func (i *Image) GetExtension() (string, error) {
	switch i.Format {
	case "png":
		return ".png", nil
	case "jpeg":
		return ".jpg", nil
	default:
		return "", xerrors.New("invalid image format")
	}
}
