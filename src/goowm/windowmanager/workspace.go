package windowmanager

import (
	"fmt"
	"goowm/config"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type Workspace struct {
	name   string
	window *xwindow.Window
}

func NewWorkspace(x *xgbutil.XUtil, conf *config.WorkspaceConfig) *Workspace {
	x.Grab()
	defer x.Ungrab()

	rg := xwindow.RootGeometry(x)
	ws, err := xwindow.Generate(x)
	if err != nil {
		panic(fmt.Errorf("Failed to generate workspace: %s", err))
	}

	err = ws.CreateChecked(x.RootWin(), 0, 20, rg.Width(), rg.Height()-20,
		xproto.CwBackPixel, 0x666666)
	if err != nil {
		panic(fmt.Errorf("Failed to create workspace: %s", err))
	}

	err = xproto.ReparentWindowChecked(x.Conn(), ws.Id, x.RootWin(), 0, 20).Check()
	if err != nil {
		panic(err)
	}

	workspace := Workspace{window: ws, name: conf.Name}

	return &workspace
}

func (ws *Workspace) Deactivate() {
	ws.window.Unmap()
}

func (ws *Workspace) Activate() {
	ws.window.Map()
}

func (ws *Workspace) WindowId() xproto.Window {
	return ws.window.Id
}
