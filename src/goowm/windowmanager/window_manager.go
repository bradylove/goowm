package windowmanager

import (
	"fmt"
	"goowm/config"
	"goowm/gwindow"
	"os"
	"os/exec"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type WindowManager struct {
	X                    *xgbutil.XUtil
	Root                 *xwindow.Window
	conf                 *config.Config
	Workspaces           []*Workspace
	ActiveWorkspaceIndex int
}

// New sets up a and returns a new *WindowManager
func New(conf *config.Config) (*WindowManager, error) {
	x, err := xgbutil.NewConnDisplay(conf.Display)
	if err != nil {
		return nil, err
	}

	wm := &WindowManager{X: x}

	for _, wc := range conf.Workspaces {
		wm.Workspaces = append(wm.Workspaces, NewWorkspace(x, wc))
	}

	fmt.Println(len(wm.Workspaces))
	wm.ActivateWorkspace(0)

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

	keybind.Initialize(x)

	xevent.MapRequestFun(wm.onMapRequest).Connect(x, x.RootWin())
	err = keybind.KeyPressFun(wm.onActivateNextWorkspace).Connect(x, x.RootWin(),
		conf.KeyBindingNextWorkspace, true)
	if err != nil {
		panic(err)
	}

	err = keybind.KeyPressFun(wm.onActivatePreviousWorkspace).Connect(x, x.RootWin(),
		conf.KeyBindingPreviousWorkspace, true)
	if err != nil {
		panic(err)
	}

	err = keybind.KeyPressFun(onShowModeLine).Connect(x, x.RootWin(), "Mod4-e", true)
	if err != nil {
		panic(err)
	}

	return wm, nil
}

func (wm *WindowManager) ActiveWorkspace() *Workspace {
	return wm.Workspaces[wm.ActiveWorkspaceIndex]
}

// ActivateWorkspace will deactivate the current active workspace if there is one
// then it will activate the workspace with the given index.
func (wm *WindowManager) ActivateWorkspace(index int) {
	wm.Workspaces[wm.ActiveWorkspaceIndex].Unmap()
	wm.Workspaces[index].Map()
	wm.ActiveWorkspaceIndex = index
}

func (wm *WindowManager) ActivateNextWorkspace() {
	var index int
	if wm.ActiveWorkspaceIndex != len(wm.Workspaces)-1 {
		index = wm.ActiveWorkspaceIndex + 1
	}

	wm.ActivateWorkspace(index)
}

func (wm *WindowManager) ActivatePreviousWorkspace() {
	index := wm.ActiveWorkspaceIndex - 1

	if index == -1 {
		index = len(wm.Workspaces) - 1
	}

	wm.ActivateWorkspace(index)
}

func (wm *WindowManager) onActivateNextWorkspace(x *xgbutil.XUtil, e xevent.KeyPressEvent) {
	wm.ActivateNextWorkspace()
}

func (wm *WindowManager) onActivatePreviousWorkspace(x *xgbutil.XUtil, e xevent.KeyPressEvent) {
	wm.ActivatePreviousWorkspace()
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

func (wm *WindowManager) onMapRequest(x *xgbutil.XUtil, e xevent.MapRequestEvent) {
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

	err = pw.CreateChecked(wm.ActiveWorkspace().WindowId(), cg.X(), cg.Y(),
		cg.Width()+12, cg.Height()+12, xproto.CwBackPixel, 0x000000)
	if err != nil {
		panic(err)
	}

	err = xproto.ReparentWindowChecked(x.Conn(), cw.Id, pw.Id, 5, 5).Check()
	if err != nil {
		panic(err)
	}

	pw.Map()
	cw.Map()

	ewmh.ActiveWindowSet(x, cw.Id)
}
