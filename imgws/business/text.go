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
	//l := len(text)
	//w := (size * 2 * l - size * 2 * l % 3 ) / 3
	w := size * 3
	//h := (size * 3 - size * 3 % 2) / 2
	h := size * 2
	println(w, h)
	pw.SetColor("none")
	mw.NewImage(uint(w), uint(h), pw)
	pw.SetColor("white")
	dw.SetFillColor(pw)
	dw.SetFont("/usr/share/fonts/default/TrueType/msyh.ttf")
	dw.SetFontSize(float64(size))
	//dw.SetTextAntialias(true)
	dw.Annotation(0, float64(size), text)
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
/*
// than "none" will remove all the transparency and replace it with the border's colour
func textEffect1() {
	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	dw := imagick.NewDrawingWand()
	defer dw.Destroy()
	pw := imagick.NewPixelWand()
	defer pw.Destroy()
	pw.SetColor("none")
	// Create a new transparent image
	mw.NewImage(350, 100, pw)
	// Set up a 72 point white font
	pw.SetColor("white")
	dw.SetFillColor(pw)
	dw.SetFont("/usr/share/fonts/default/TrueType/verdana.ttf")
	dw.SetFontSize(72)
	// Add a black outline to the text
	pw.SetColor("black")
	dw.SetStrokeColor(pw)
	// Turn antialias on - not sure this makes a difference
	dw.SetTextAntialias(true)
	// Now draw the text
	dw.Annotation(25, 65, "Magick")
	// Draw the image on to the mw
	mw.DrawImage(dw)
	// Trim the image down to include only the text
	mw.TrimImage(0)
	// equivalent to the command line +repage
	mw.ResetImagePage("")
	// Make a copy of the text image
	cw := mw.Clone()
	// Set the background colour to blue for the shadow
	pw.SetColor("blue")
	mw.SetImageBackgroundColor(pw)
	// Opacity is a real number indicating (apparently) percentage
	mw.ShadowImage(70, 4, 5, 5)
	// Composite the text on top of the shadow
	mw.CompositeImage(cw, imagick.COMPOSITE_OP_OVER, 5, 5)
	cw.Destroy()
	cw = imagick.NewMagickWand()
	defer cw.Destroy()
	// Create a new image the same size as the text image and put a solid colour
	// as its background
	pw.SetColor("rgb(125,215,255)")
	cw.NewImage(mw.GetImageWidth(), mw.GetImageHeight(), pw)
	// Now composite the shadowed text over the plain background
	cw.CompositeImage(mw, imagick.COMPOSITE_OP_OVER, 0, 0)
	// and write the result
	*/
