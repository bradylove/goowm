package panel

import (
	"bytes"
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
	"github.com/martine/gocairo/cairo"
)

type Panel struct {
	x      *xgbutil.XUtil
	window *xwindow.Window
	conf   *config.Config
}

func New(x *xgbutil.XUtil, names []string) (*Panel, error) {
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

	var xPos, yPos int
	for _, n := range names {
		width := renderWorkspace(x, win, n, xPos, yPos)
		xPos += width
	}

	return &Panel{x: x, window: win}, nil
}

func (p *Panel) Run() {
	p.window.Map()
}

func renderWorkspace(x *xgbutil.XUtil, p *xwindow.Window, n string, xPos, yPos int) int {
	buf, err := renderText(n)
	if err != nil {
		panic(err)
	}

	// render.Text(x, p.Id, currentTime, font, size, wg.Width()-w-3, 0)
	img, err := xgraphics.NewBytes(x, buf.Bytes())
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

	buf, err := renderText(currentTime)
	if err != nil {
		panic(err)
	}

	// render.Text(x, p.Id, currentTime, font, size, wg.Width()-w-3, 0)
	img, err := xgraphics.NewBytes(x, buf.Bytes())
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

func renderText(text string) (*bytes.Buffer, error) {
	surf := cairo.ImageSurfaceCreate(cairo.FormatARGB32, 60, 30)
	
	cr := cairo.Create(surf.Surface)

	cr.SetSourceRGB(0.2, 0.2, 0.2)
	// cr.SetSourceRGB(0.9, 0.9, 0.9)
	cr.PaintWithAlpha(1.0)

	cr.SetAntialias(cairo.AntialiasBest)

	cr.SetSourceRGB(1, 1, 1)
	cr.SelectFontFace("SourceCodePro-Bold", cairo.FontSlantNormal, cairo.FontWeightNormal)
	cr.SetFontSize(14)
	cr.MoveTo(60/10, 30/2)
	cr.ShowText(text)

	buf := bytes.NewBuffer(nil)
	if err := surf.WriteToPNG(buf); err != nil {
		return nil, err
	}

	return buf, nil
}
