package windowmanager

import (
	"goowm/config"
	"os"
	"os/exec"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type WindowManager struct {
	X                *xgbutil.XUtil
	Root             *xwindow.Window
	conf             *config.Config
	WorkspaceManager *WorkspaceManager
}

// New sets up a and returns a new *WindowManager
func New(conf *config.Config) (*WindowManager, error) {
	x, err := xgbutil.NewConnDisplay(conf.Display)
	if err != nil {
		return nil, err
	}

	wm := &WindowManager{
		X:                x,
		WorkspaceManager: NewWorkspaceManager(len(conf.Workspaces)),
	}

	wm.WorkspaceManager.Add(x, conf.Workspaces...)
	wm.WorkspaceManager.Activate(0)

	root := xwindow.New(x, x.RootWin())
	panel, err := NewPanel(x, wm)
	if err != nil {
		panic(err)
	}
	panel.Run()

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
	xevent.ConfigureRequestFun(wm.onConfigureRequest).Connect(x, x.RootWin())

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

	err = keybind.KeyPressFun(onExecLauncher).Connect(x, x.RootWin(), "Mod4-e", true)
	if err != nil {
		panic(err)
	}

	return wm, nil
}

func (wm *WindowManager) activeWorkspace() *Workspace {
	return wm.WorkspaceManager.workspaces[wm.WorkspaceManager.activeIndex]
}

// func (wm *WindowManager) activateWorkspace(index int) {
// 	wm.WorkspaceManager.Activate(index)
// }

func (wm *WindowManager) activateNextWorkspace() {
	wm.WorkspaceManager.Activate(wm.WorkspaceManager.NextIndex())
}

func (wm *WindowManager) activatePreviousWorkspace() {
	wm.WorkspaceManager.Activate(wm.WorkspaceManager.PreviousIndex())
}

func (wm *WindowManager) onActivateNextWorkspace(x *xgbutil.XUtil, e xevent.KeyPressEvent) {
	wm.activateNextWorkspace()
}

func (wm *WindowManager) onActivatePreviousWorkspace(x *xgbutil.XUtil, e xevent.KeyPressEvent) {
	wm.activatePreviousWorkspace()
}

func (wm *WindowManager) Run() {
	xevent.Main(wm.X)
}

func onExecLauncher(x *xgbutil.XUtil, e xevent.KeyPressEvent) {
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

	w := NewWindow(x, wm.activeWorkspace().WindowId(), e.Window)
	w.Draw()
}

func (wm *WindowManager) onConfigureRequest(x *xgbutil.XUtil, ev xevent.ConfigureRequestEvent) {
	xwindow.New(x, ev.Window).Configure(int(ev.ValueMask), int(ev.X), int(ev.Y),
		int(ev.Width), int(ev.Height), ev.Sibling, ev.StackMode)
}
