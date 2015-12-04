package img4g

/*
#cgo CFLAGS: -std=c99
#cgo CPPFLAGS: -I/usr/local/include/GraphicsMagick
#cgo LDFLAGS: -L/usr/local/lib -L/usr/local/lib64 -L/usr/lib -L/usr/lib64 -lGraphicsMagickWand -lGraphicsMagick -lfreetype -ljpeg -lpng12  -lbz2 -lxml2  -lz -lm -lgomp  -lpthread -Wl,-no_compact_unwind
#include <wand/magick_wand.h>
#include "cmagick.h"
*/
import "C"
import (
	"errors"
	"image"
	"image/color"
	"image/png"

	cat "github.com/ctripcorp/cat.go"
)

var (
	ErrIllegalFormat = errors.New("illegal format")
)

func (this *Image) Write(p []byte) (int, error) {
	this.Blob = append(this.Blob, p...)
	return len(this.Blob), nil
}

func NewImage(width, height int, format string, CAT cat.Cat) (*Image, error) {
	switch format {
	case "PNG":
		return NewImageAsPNG(width, height, CAT)
	default:
		return nil, ErrIllegalFormat
	}
}

func NewImageAsPNG(width, height int, CAT cat.Cat) (*Image, error) {
	i := &Image{
		Format: "PNG",
		Cat:    CAT,
	}
	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if x == y {
				rgba.Set(x, y, color.RGBA{0, 0, 0, 255})
			} else {
				rgba.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}
	err := png.Encode(i, rgba)
	if err != nil {
		return nil, err
	}
	return i, nil
}
