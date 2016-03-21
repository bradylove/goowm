package main

import (
	"fmt"
	"goowm/render"
	"os/exec"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/motif"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
)

const (
	LauncherHeight          = 40
	LauncherWidth           = 300
	LauncherBackgroundColor = 0xffffff
	LauncherBorderColor     = 0x333333
)

func main() {
	x, err := xgbutil.NewConnDisplay(":0")
	if err != nil {
		panic(err)
	}

	keybind.Initialize(x)

	rg := xwindow.RootGeometry(x)

	xPos := (rg.Width() / 2) - (LauncherWidth / 2)
	yPos := (rg.Height() / 2) - (LauncherHeight / 2)

	parent, err := xwindow.Generate(x)
	if err != nil {
		panic(err)
	}

	input, err := xwindow.Generate(x)
	if err != nil {
		panic(err)
	}

	parent.CreateChecked(x.RootWin(), xPos, yPos, LauncherWidth, LauncherHeight,
		xproto.CwBackPixel, LauncherBorderColor)

	input.CreateChecked(parent.Id, 3, 3, LauncherWidth-6, LauncherHeight-6,
		xproto.CwBackPixel, LauncherBackgroundColor)

	motif.WmHintsSet(x, parent.Id, &motif.Hints{
		Flags:      motif.HintDecorations,
		Decoration: motif.DecorationNone,
	})

	if err := input.Listen(xproto.EventMaskKeyPress); err != nil {
		panic(err)
	}

	tc := render.NewColor(0, 0, 0)
	bc := render.NewColor(255, 255, 255)
	var textWin *xwindow.Window
	var drawInputedChars = func(txt string) {
		oldWin := textWin

		imgData, w, h, err := BuildTextImg(txt, "SourceCodePro", 18.0, tc, bc)
		if err != nil {
			panic(err)
		}

		fmt.Println("Width: ", w)
		fmt.Println("Height:", h)

		img, err := xgraphics.NewBytes(x, imgData)
		if err != nil {
			panic(err)
		}

		textWin, err = xwindow.Generate(x)
		if err != nil {
			panic(err)
		}

		if err := textWin.CreateChecked(input.Id, 0, 0, 1, 1, xproto.CwBackPixel, 0x333333); err != nil {
			panic(err)
		}

		textWin.MoveResize(5, 7, w, 30)

		img.XSurfaceSet(textWin.Id)
		img.XDraw()
		img.XPaint(textWin.Id)
		img.Destroy()

		textWin.Map()
		if oldWin != nil {
			oldWin.Destroy()
		}
		fmt.Println("You should see text...")
	}

	var txt string
	xevent.KeyPressFun(
		func(x *xgbutil.XUtil, e xevent.KeyPressEvent) {
			if keybind.KeyMatch(x, "Escape", e.State, e.Detail) {
				fmt.Println("Exiting...")
				xevent.Quit(x)
				return
			}

			if keybind.KeyMatch(x, "BackSpace", e.State, e.Detail) && len(txt) > 0 {
				txt = txt[:len(txt)-1]
				drawInputedChars(txt)
				fmt.Println(txt)

				return
			}

			if keybind.KeyMatch(x, "Return", e.State, e.Detail) && len(txt) > 0 {
				launchExecutable(txt)
				fmt.Println("Exiting...")
				xevent.Quit(x)
			}

			key := keybind.LookupString(x, e.State, e.Detail)
			fmt.Println(key)

			if len(key) > 1 {
				return
			}

			r, _ := utf8.DecodeRuneInString(key[0:])

			if unicode.IsPrint(r) {
				txt += key
				drawInputedChars(txt)
				fmt.Println(txt)
			}
		}).Connect(x, input.Id)

	input.Map()
	parent.Map()

	if err := ewmh.ActiveWindowReq(x, input.Id); err != nil {
		panic(err)
	}

	time.Sleep(100 * time.Millisecond)
	if err := keybind.GrabKeyboard(x, input.Id); err != nil {
		panic(err)
	}
	defer keybind.GrabKeyboard(x, input.Id)

	xevent.Main(x)
}

func launchExecutable(name string) {
	path, err := exec.LookPath(name)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(path)
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
}
