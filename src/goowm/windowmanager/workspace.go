package windowmanager

import (
	"fmt"
	"goowm/config"
	"image"
	"io/ioutil"

	"github.com/BurntSushi/freetype-go/freetype/truetype"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xgraphics"
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

	err = ws.CreateChecked(x.RootWin(), 0, 0, rg.Width(), rg.Height(),
		xproto.CwBackPixel, 0xffffff)
	if err != nil {
		panic(fmt.Errorf("Failed to create workspace: %s", err))
	}

	err = xproto.ReparentWindowChecked(x.Conn(), ws.Id, x.RootWin(), 0, 0).Check()
	if err != nil {
		panic(err)
	}

	workspace := Workspace{window: ws, name: conf.Name}
	workspace.renderName(x)

	return &workspace
}

func (ws *Workspace) Unmap() {
	ws.window.Unmap()
}

func (ws *Workspace) Map() {
	ws.window.Map()
}

func (ws *Workspace) WindowId() xproto.Window {
	return ws.window.Id
}

// TODO: Move everything below this line somewhere else, it doesn't belong here
type Color struct {
	R uint32
	G uint32
	B uint32
	A uint32
}

func (c Color) RGBA() (r, g, b, a uint32) {
	return c.R, c.G, c.B, c.A
}

func (ws *Workspace) renderName(x *xgbutil.XUtil) error {
	win, err := xwindow.Generate(x)
	if err != nil {
		panic(fmt.Errorf("Failed to generate workspace: %s", err))
	}

	err = win.CreateChecked(ws.window.Id, 0, 0, 1, 1, xproto.CwBackPixel, 0x000000)
	if err != nil {
		panic(fmt.Errorf("Failed to create workspace: %s", err))
	}

	fd, err := ioutil.ReadFile("resources/DejaVuSans.ttf")
	if err != nil {
		return err
	}

	font, err := truetype.Parse(fd)
	if err != nil {
		return err
	}

	size := 28.0

	ew, eh := xgraphics.Extents(font, size, ws.name)
	img := xgraphics.New(x, image.Rect(0, 0, ew, eh))
	xgraphics.BlendBgColor(img, Color{R: 255, G: 255, B: 255})

	_, _, err = img.Text(0, 0, Color{R: 100, G: 100, B: 100}, size, font, ws.name)
	if err != nil {
		return err
	}

	wg, err := ws.window.Geometry()
	if err != nil {
		return err
	}

	win.MoveResize(wg.Width()-ew-5, wg.Height()-eh-5, ew, eh)

	img.XSurfaceSet(win.Id)
	img.XDraw()
	img.XPaint(win.Id)
	img.Destroy()

	win.Map()

	return nil
}
