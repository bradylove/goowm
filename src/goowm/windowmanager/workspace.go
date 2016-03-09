package windowmanager

import (
	"fmt"
	"goowm/config"
	"goowm/render"
	"io/ioutil"

	"github.com/BurntSushi/freetype-go/freetype/truetype"
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
	workspace.renderName(x)

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

func (ws *Workspace) renderName(x *xgbutil.XUtil) error {
	size := 28.0

	fd, err := ioutil.ReadFile("/usr/share/fonts/TTF/SourceCodePro-Medium.ttf")
	if err != nil {
		panic(err)
	}

	font, err := truetype.Parse(fd)
	if err != nil {
		panic(err)
	}

	w, h := render.Extents(ws.name, font, size)
	wg, err := ws.window.Geometry()
	if err != nil {
		return err
	}

	render.Text(x, ws.WindowId(), ws.name, font, size, wg.Width()-w-5, wg.Height()-h-5)
	return nil
}
