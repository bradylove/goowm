package gwindow

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type Window struct {
	*xwindow.Window
	X *xgbutil.XUtil
}

// New retreives the xwindow from X and returns a new instance of *Window
func New(x *xgbutil.XUtil, id xproto.Window) *Window {
	xwin := xwindow.New(x, id)

	return &Window{
		Window: xwin,
		X:      x,
	}
}

func (w *Window) Maximize() error {
	rg := xwindow.RootGeometry(w.X)
	w.MoveResize(rg.X(), rg.Y(), rg.Width(), rg.Height())
	return nil
}
