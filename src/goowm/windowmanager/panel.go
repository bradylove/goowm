package windowmanager

import (
	"fmt"
	"goowm/broadcaster"
	"goowm/config"
	"goowm/render"
	"io/ioutil"
	"log"
	"time"

	"github.com/BurntSushi/freetype-go/freetype/truetype"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type Panel struct {
	x      *xgbutil.XUtil
	window *xwindow.Window
	conf   *config.Config
	wm     *WindowManager
}

func NewPanel(x *xgbutil.XUtil, wm *WindowManager) (*Panel, error) {
	log.Println("Creating new panel")

	rg := xwindow.RootGeometry(x)
	win, err := xwindow.Generate(x)
	if err != nil {
		return nil, err
	}

	err = win.CreateChecked(x.RootWin(), 0, 0, rg.Width(), 20, xproto.CwBackPixel, 0x333333)
	if err != nil {
		return nil, err
	}

	runClock(x, win)

	p := &Panel{x: x, window: win, wm: wm}
	p.drawWorkspaces()

	broadcaster.Listen(broadcaster.EventWorkspaceChanged, func() {
		fmt.Println("Workspaces Changed")
		p.drawWorkspaces()
	})

	return p, nil
}

func (p *Panel) drawWorkspaces() {
	var xPos, yPos int
	for i, w := range p.wm.WorkspaceManager.Workspaces {
		tc := render.NewColor(230, 230, 230)

		if i == p.wm.WorkspaceManager.ActiveIndex() {
			tc = render.NewColor(255, 153, 0)
		}

		width := renderWorkspace(p.x, p.window, w.Name(), xPos, yPos, tc)
		xPos += width
	}
}

func (p *Panel) Run() {
	p.window.Map()
}

func renderWorkspace(x *xgbutil.XUtil, p *xwindow.Window, n string, xPos, yPos int, tc render.Color) int {
	textImg, err := render.Text(n, "SourceCodePro", tc, 14.0, 60, 30)
	if err != nil {
		panic(err)
	}

	img, err := xgraphics.NewBytes(x, textImg)
	if err != nil {
		panic(err)
	}

	win, err := xwindow.Generate(x)
	if err != nil {
		panic(err)
	}

	err = win.CreateChecked(p.Id, 0, 0, 1, 1, xproto.CwBackPixel, 0xffffff)
	if err != nil {
		panic(err)
	}

	win.MoveResize(xPos, 0, 60, 30)

	img.XSurfaceSet(win.Id)
	img.XDraw()
	img.XPaint(win.Id)
	img.Destroy()

	win.Map()
	return 60
}

func runClock(x *xgbutil.XUtil, p *xwindow.Window) {
	ticker := time.NewTicker(1 * time.Second)

	drawClock(x, p)

	go func() {
		for _ = range ticker.C {
			drawClock(x, p)
		}
	}()
}

func drawClock(x *xgbutil.XUtil, p *xwindow.Window) {
	log.Println("Updating clock")

	currentTime := time.Now().Format("15:04")
	size := 18.0

	fd, err := ioutil.ReadFile("/usr/share/fonts/TTF/SourceCodePro-Medium.ttf")
	if err != nil {
		panic(err)
	}

	font, err := truetype.Parse(fd)
	if err != nil {
		panic(err)
	}

	w, _ := render.Extents(currentTime, font, size)
	wg, err := p.Geometry()
	if err != nil {
		panic(err)
	}

	tc := render.NewColor(230, 230, 230)
	textImg, err := render.Text(currentTime, "SourceCodePro", tc, 14.0, 60, 30)
	if err != nil {
		panic(err)
	}

	img, err := xgraphics.NewBytes(x, textImg)
	if err != nil {
		panic(err)
	}

	win, err := xwindow.Generate(x)
	if err != nil {
		panic(err)
	}

	err = win.CreateChecked(p.Id, 0, 0, 1, 1, xproto.CwBackPixel, 0xffffff)
	if err != nil {
		panic(err)
	}

	win.MoveResize(wg.Width()-w-3, 0, 60, 30)

	img.XSurfaceSet(win.Id)
	img.XDraw()
	img.XPaint(win.Id)
	img.Destroy()

	win.Map()
}
