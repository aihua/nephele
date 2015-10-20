package business

import (
	"github.com/gographics/imagick/imagick"
)

func getTextImage(text string) []byte {
	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	dw := imagick.NewDrawingWand()
	defer dw.Destroy()
	pw := imagick.NewPixelWand()
	defer pw.Destroy()
	pw.SetColor("none")
	mw.NewImage(320, 100, pw)
	dw.SetFont("/usr/share/fonts/default/TrueType/verdana.ttf")
	dw.SetFontSize(72)
	dw.Annotation(25, 65, "Magick")
	mw.DrawImage(dw)
	mw.WriteImage("text_plain.png")
	return mw.GetImageBlob()
}
