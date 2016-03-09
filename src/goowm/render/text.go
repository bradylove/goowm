package render

import (
	"fmt"
	"image"

	"github.com/BurntSushi/freetype-go/freetype/truetype"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
)

func Extents(text string, font *truetype.Font, size float64) (int, int) {
	return xgraphics.Extents(font, size, text)
}

// Redo this using cairo
func Text(x *xgbutil.XUtil, parentId xproto.Window, text string, font *truetype.Font,
	fontSize float64, xPos, yPos int) error {

	win, err := xwindow.Generate(x)
	if err != nil {
		return fmt.Errorf("Failed to generate text window: %s", err)
	}

	err = win.CreateChecked(parentId, 0, 0, 1, 1, xproto.CwBackPixel, 0x666666)
	if err != nil {
		return fmt.Errorf("Failed to create text window: %s", err)
	}

	ew, eh := Extents(text, font, fontSize)

	img := xgraphics.New(x, image.Rect(0, 0, ew, eh))
	xgraphics.BlendBgColor(img, Color{R: 255, G: 255, B: 255})

	_, _, err = img.Text(0, 0, Color{R: 100, G: 100, B: 100}, fontSize, font, text)
	if err != nil {
		return err
	}

	win.MoveResize(xPos, yPos, ew, eh)

	img.XSurfaceSet(win.Id)
	img.XDraw()
	img.XPaint(win.Id)
	img.Destroy()

	win.Map()

	return nil
}
