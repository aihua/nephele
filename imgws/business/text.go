package business

import (
	"github.com/gographics/imagick/imagick"
)

func getTextImage(text string, size int) []byte {
	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	dw := imagick.NewDrawingWand()
	defer dw.Destroy()
	pw := imagick.NewPixelWand()
	defer pw.Destroy()
	l := len(text)
	w := (size * 2 * l - size * 2 * l % 3 ) / 3
	h := (size * 3 - size * 3 % 2) / 2
	println(w, h)
	pw.SetColor("none")
	mw.NewImage(uint(w), uint(h), pw)
	dw.SetFont("/usr/share/fonts/default/TrueType/verdana.ttf")
	dw.SetFontSize(float64(size))
	dw.Annotation(0, float64(size), text)
	mw.DrawImage(dw)
	mw.SetImageFormat("PNG")
	return mw.GetImageBlob()
}
