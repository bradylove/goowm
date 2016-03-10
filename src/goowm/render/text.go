package render

import (
	"bytes"

	"github.com/BurntSushi/freetype-go/freetype/truetype"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/martine/gocairo/cairo"
)

func Extents(text string, font *truetype.Font, size float64) (int, int) {
	return xgraphics.Extents(font, size, text)
}

func Text(text, font string, color Color, size float64, width, height int) ([]byte, error) {
	surf := cairo.ImageSurfaceCreate(cairo.FormatARGB32, width, height)

	cr := cairo.Create(surf.Surface)

	fr, fg, fb := color.ToFloat64s()
	cr.SetSourceRGB(0.2, 0.2, 0.2)
	cr.PaintWithAlpha(1.0)
	cr.SetAntialias(cairo.AntialiasBest)
	cr.SetSourceRGB(fr, fg, fb)
	cr.SelectFontFace(font, cairo.FontSlantNormal, cairo.FontWeightNormal)
	cr.SetFontSize(size)
	cr.MoveTo(float64(width)/10, float64(height)/2)
	cr.ShowText(text)

	buf := bytes.NewBuffer(nil)
	if err := surf.WriteToPNG(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
