package main

import (
	"bytes"
	"goowm/render"

	"github.com/martine/gocairo/cairo"
)

func BuildTextImg(txt, font string, size float64, txtColor, bgColor render.Color) ([]byte, int, int, error) {
	surf := cairo.ImageSurfaceCreate(cairo.FormatARGB32, 300, 30)
	cr := cairo.Create(surf.Surface)

	var ext cairo.TextExtents
	cr.SelectFontFace(font, cairo.FontSlantNormal, cairo.FontWeightNormal)
	cr.SetFontSize(size)
	cr.TextExtents(txt, &ext)

	surf = cairo.ImageSurfaceCreate(cairo.FormatARGB32, 300, 30)
	cr = cairo.Create(surf.Surface)

	fr, fg, fb := bgColor.ToFloat64s()
	cr.SetSourceRGB(fr, fg, fb)
	cr.PaintWithAlpha(1.0)

	fr, fg, fb = txtColor.ToFloat64s()
	cr.SetSourceRGB(fr, fg, fb)
	cr.SelectFontFace(font, cairo.FontSlantNormal, cairo.FontWeightNormal)
	cr.SetFontSize(size)
	cr.MoveTo(0, 30/2)
	cr.ShowText(txt)

	buf := bytes.NewBuffer(nil)
	if err := surf.WriteToPNG(buf); err != nil {
		return nil, 0, 0, err
	}

	return buf.Bytes(), int(ext.Width) + 2, int(ext.Height) + 2, nil
}
