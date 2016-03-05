package windowmanager

import (
	"goowm/config"
	"goowm/gwindow"
	"os"
	"os/exec"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type WindowManager struct {
	X    *xgbutil.XUtil
	Root *xwindow.Window
	conf *config.Config
}

// New sets up a and returns a new *WindowManager
func New(conf *config.Config) (*WindowManager, error) {
	x, err := xgbutil.NewConnDisplay(conf.Display)
	if err != nil {
		return nil, err
	}

	root := xwindow.New(x, x.RootWin())
	evMasks := xproto.EventMaskPropertyChange |
		xproto.EventMaskFocusChange |
		xproto.EventMaskButtonPress |
		xproto.EventMaskButtonRelease |
		xproto.EventMaskStructureNotify |
		xproto.EventMaskSubstructureRedirect |
		xproto.EventMaskSubstructureRedirect

	if err := root.Listen(evMasks); err != nil {
		panic(err)
	}

	mousebind.Initialize(x)
	keybind.Initialize(x)

	xevent.MapRequestFun(onMapRequest).Connect(x, x.RootWin())
	xevent.ConfigureRequestFun(onConfigureRequest).Connect(x, x.RootWin())

	err = keybind.KeyPressFun(onShowModeLine).Connect(x, x.RootWin(), "Mod4-x", true)
	if err != nil {
		panic(err)
	}

	mousebind.ButtonPressFun(
		func(x *xgbutil.XUtil, e xevent.ButtonPressEvent) {
			if e.Child != 0 {
				gwindow.New(x, e.Child).Maximize()
			}
		}).Connect(x, x.RootWin(), "Mod4-1", false, true)

	return &WindowManager{
		X:    x,
		Root: root,
	}, nil
}

func (wm *WindowManager) Run() {
	xevent.Main(wm.X)
}

func onShowModeLine(x *xgbutil.XUtil, e xevent.KeyPressEvent) {
	cmd := exec.Command("dmenu_run", "-b")

	env := os.Environ()
	env = append(env, "DISPLAY=:1")

	cmd.Env = env

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func onMapRequest(x *xgbutil.XUtil, e xevent.MapRequestEvent) {
	x.Grab()
	defer x.Ungrab()

	cw := gwindow.New(x, e.Window)
	cg, err := cw.Geometry()
	if err != nil {
		panic(err)
	}

	pw, err := xwindow.Generate(x)
	if err != nil {
		panic(err)
	}

	err = pw.CreateChecked(x.RootWin(), cg.X(), cg.Y(),
		cg.Width()+12, cg.Height()+12, xproto.CwBackPixel, 0x66ff33)
	if err != nil {
		panic(err)
	}

	_ = xproto.ReparentWindowChecked(x.Conn(), cw.Id, pw.Id, 5, 5)
	if err != nil {
		panic(err)
	}

	cw.Map()
	pw.Map()
}

func onConfigureRequest(x *xgbutil.XUtil, ev xevent.ConfigureRequestEvent) {
	xwindow.New(x, ev.Window).Configure(int(ev.ValueMask), int(ev.X), int(ev.Y),
		int(ev.Width), int(ev.Height), ev.Sibling, ev.StackMode)
}
