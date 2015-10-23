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
	w := size * 3 * (l+6)
	h := size * 2
	println(l, w ,h)
	pw.SetColor("none")
	mw.NewImage(uint(w), uint(h), pw)
	pw.SetColor("white")
	dw.SetFillColor(pw)
	dw.SetFont("/usr/share/fonts/default/TrueType/msyh.ttf")
	dw.SetFontSize(float64(size))
	dw.Annotation(0, float64(size), "ctrip © " + text)
	mw.DrawImage(dw)
	mw.TrimImage(0)

	mw.ResetImagePage("")
	cw := mw.Clone()
	pw.SetColor("black")
	mw.SetImageBackgroundColor(pw)
	mw.ShadowImage(100, 1, 0, 0)
	mw.CompositeImage(cw, imagick.COMPOSITE_OP_OVER, 1, 1)
	cw.Destroy()

	mw.SetImageFormat("PNG")
	return mw.GetImageBlob()
}
