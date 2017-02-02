package image

import (
	"image"
	"os"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/riff"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type Image struct {
	Path   string
	Type   string
	Width  int64
	Height int64
	Size   int64
}

func Stat(path string) (*Image, error) {

	img := &Image{Path: path}

	f, err := os.Stat(path)
	if err != nil {
		return img, err
	}
	img.Size = f.Size()

	fh, err := os.Open(path)
	if err != nil {
		return img, err
	}

	c, t, err := image.DecodeConfig(fh)
	if err != nil {
		return img, err
	}

	img.Type = t
	img.Width = int64(c.Width)
	img.Height = int64(c.Height)
	return img, nil
}
