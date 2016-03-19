package windowmanager

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type Window struct {
	X           *xgbutil.XUtil
	child       *xwindow.Window
	parent      *xwindow.Window
	workspaceId xproto.Window
}

// New retreives the xwindow from X and returns a new instance of *Window
func NewWindow(x *xgbutil.XUtil, workspaceId xproto.Window, id xproto.Window) *Window {
	child := xwindow.New(x, id)
	pw, err := xwindow.Generate(x)
	if err != nil {
		panic(err)
	}

	win := Window{
		X:           x,
		child:       child,
		parent:      pw,
		workspaceId: workspaceId,
	}

	win.bindEvents()
	win.reparentChild()
	win.Activate()

	return &win
}

func (w *Window) reparentChild() error {
	cg, err := w.child.Geometry()
	if err != nil {
		panic(err)
	}

	err = w.parent.CreateChecked(w.workspaceId, cg.X()+100, cg.Y()+100,
		cg.Width()+8, cg.Height()+8, xproto.CwBackPixel, 0x111111)
	if err != nil {
		return err
	}

	err = xproto.ReparentWindowChecked(w.X.Conn(), w.child.Id, w.parent.Id, 4, 4).Check()
	if err != nil {
		return err
	}

	return nil
}

func (w *Window) bindEvents() {
	evMask := xproto.EventMaskStructureNotify

	if err := w.child.Listen(evMask); err != nil {
		panic(err)
	}

	xevent.DestroyNotifyFun(func(x *xgbutil.XUtil, e xevent.DestroyNotifyEvent) {
		w.Destroy()
	}).Connect(w.X, w.child.Id)
}

func (w *Window) Maximize() error {
	rg := xwindow.RootGeometry(w.X)
	w.parent.MoveResize(rg.X(), rg.Y(), rg.Width(), rg.Height())
	return nil
}

func (w *Window) Destroy() {
	w.parent.Destroy()
}

func (w *Window) Draw() {
	w.parent.Map()
	w.child.Map()
}

func (w *Window) Activate() error {
	return ewmh.ActiveWindowSet(w.X, w.child.Id)
}
