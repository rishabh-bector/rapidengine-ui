package main

//   --------------------------------------------------
//   Main UI for Rapid Engine
//   --------------------------------------------------

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"rapidengine/cmd"
	"rapidengine/input"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"nuklear-golang/nk"
)

var engine *cmd.Engine

// Windows
var currentWindow = 1

// Panel sizes
var panelSize = float32(0.2)
var sideBarSize = float32(0.025)

// idk
var sliderVal = float32(50)

const (
	winWidth  = 3840
	winHeight = 2160

	maxVertexBuffer  = 512 * 1024
	maxElementBuffer = 128 * 1024
)

func init() {
	runtime.LockOSThread()
}

var ctx *nk.Context
var state *State

var avenirFont *nk.Font
var aromaFont *nk.Font

func main() {
	ScreenWidth := 3840
	ScreenHeight := 2160
	config := cmd.NewEngineConfig(ScreenWidth, ScreenHeight, 3)
	config.FullScreen = false
	config.VSync = true
	config.AntiAliasing = true
	config.PolygonLines = false
	config.GammaCorrection = true
	config.ShowFPS = false
	engine = cmd.NewEngine(&config, render)
	engine.Renderer.MainCamera.SetSpeed(60.0 / 20.0)
	scene := engine.SceneControl.NewScene("main")

	gl.Init()

	ctx = nk.NkPlatformInit(engine.Renderer.Window, nk.PlatformInstallCallbacks)

	avenirFile, _ := os.Open("avenir-next-regular.ttf")
	avenirBytes, _ := ioutil.ReadAll(avenirFile)

	aromaFile, _ := os.Open("aroma-bold.ttf")
	aromaBytes, _ := ioutil.ReadAll(aromaFile)

	atlas := nk.NewFontAtlas()
	nk.NkFontStashBegin(&atlas)

	avenirFont = nk.NkFontAtlasAddFromBytes(atlas, avenirBytes, 40, nil)
	aromaFont = nk.NkFontAtlasAddFromBytes(atlas, aromaBytes, 40, nil)

	nk.NkFontStashEnd()
	if avenirFont != nil {
		nk.NkStyleSetFont(ctx, avenirFont.Handle())
	}

	state = &State{
		bgColor: nk.NkRgba(28, 48, 62, 255),
	}

	engine.SceneControl.InstanceScene(scene)
	engine.SceneControl.SetCurrentScene(scene)
	engine.EnableLighting()

	engine.Initialize()
	engine.StartRenderer()
	<-engine.Done()
}

func render(renderer *cmd.Renderer, inputs *input.Input) {
	if inputs.RightMouseButton {
		renderer.MainCamera.DefaultControls(inputs)
	}
	gfxMain(renderer.Window, ctx, state)
}

func gfxMain(win *glfw.Window, ctx *nk.Context, state *State) {
	nk.NkPlatformNewFrame()

	//   --------------------------------------------------
	//   Left sidebar
	//   --------------------------------------------------

	bounds := nk.NkRect(0, 0, sideBarSize*winWidth, winHeight)
	nk.NkBegin(ctx, "Bar", bounds, nk.WindowNoScrollbar)

	nk.NkEnd(ctx)

	//   --------------------------------------------------
	//   Left panel
	//   --------------------------------------------------

	bounds = nk.NkRect(sideBarSize*winWidth, 0, panelSize*winWidth, winHeight)
	nk.NkBegin(ctx, "Demo", bounds, nk.WindowNoScrollbar)

	//cb := nk.NkWindowGetCanvas(ctx)
	//nk.NkStrokeLine(cb, sideBarSize*winWidth+10, 0, sideBarSize*winWidth+10, winHeight, 5, nk.NkRgb(255, 255, 255))
	//nk.NkStrokeLine(cb, sideBarSize*winWidth+10+panelSize*winWidth-sliderVal, 0, sideBarSize*winWidth+10+panelSize*winWidth-sliderVal, winHeight, 5, nk.NkRgb(255, 255, 255))

	if currentWindow == 2 {
		leftMaterial()
	}

	nk.NkEnd(ctx)

	//   --------------------------------------------------
	//   Right panel
	//   --------------------------------------------------

	bounds = nk.NkRect(winWidth-(panelSize*winWidth), 0, panelSize*winWidth, winHeight)
	nk.NkBegin(ctx, "", bounds, nk.WindowNoScrollbar)

	nk.NkLayoutRowDynamic(ctx, 100, 1)
	nk.NkLabel(ctx, "Components", nk.TextAlignCentered|nk.TextAlignMiddle)

	ratio := []float32{0.2, 0.8}
	nk.NkLayoutRow(ctx, nk.Dynamic, 50, 2, ratio)
	nk.NkLabel(ctx, fmt.Sprintf("%v", sliderVal), nk.TextAlignCentered|nk.TextAlignMiddle)
	nk.NkSliderFloat(ctx, 0, &sliderVal, 100, 1)

	nk.NkLayoutRowDynamic(ctx, 50, 1)
	nk.NkButtonLabel(ctx, "click me")

	if currentWindow == 2 {
		rightMaterial()
	}

	nk.NkEnd(ctx)

	//   --------------------------------------------------
	//   Top panel
	//   --------------------------------------------------

	bounds = nk.NkRect(panelSize*winWidth+sideBarSize*winWidth, 0, winWidth-(2*panelSize*winWidth)-(sideBarSize*winWidth), 100)
	nk.NkBegin(ctx, "engine", bounds, nk.WindowNoScrollbar|nk.WindowNoScrollbar)

	ratio = []float32{0.2, 0.2, 0.2, 0.2, 0.2}

	cb := nk.NkWindowGetCanvas(ctx)

	// Selector lines
	ind := currentWindow
	if ind > 2 {
		ind += 1
	}
	lineWidth := (winWidth - (2 * panelSize * winWidth) - (sideBarSize * winWidth)) / 5
	nk.NkStrokeLine(
		cb,
		panelSize*winWidth+sideBarSize*winWidth+lineWidth*float32(ind-1),
		96,
		panelSize*winWidth+sideBarSize*winWidth+lineWidth*float32(ind),
		96, 5, nk.NkRgb(255, 255, 255),
	)
	nk.NkStrokeLine(
		cb,
		panelSize*winWidth+sideBarSize*winWidth+lineWidth*float32(ind-1),
		2,
		panelSize*winWidth+sideBarSize*winWidth+lineWidth*float32(ind),
		2, 5, nk.NkRgb(255, 255, 255),
	)

	//nk.NkLayoutRowDynamic(ctx, 50, 1)
	//nk.NkGroupBegin(ctx, "E", nk.WindowBorder)
	//nk.NkGroupEnd(ctx)

	nk.NkLayoutRow(ctx, nk.Dynamic, 80, 5, ratio)
	if nk.NkButtonLabel(ctx, "Scene") == 1 {
		currentWindow = 1
	}
	if nk.NkButtonLabel(ctx, "Materials") == 1 {
		currentWindow = 2
	}

	nk.NkStyleSetFont(ctx, aromaFont.Handle())
	nk.NkLabel(ctx, "RAPID ENGINE", nk.TextAlignCentered|nk.TextAlignMiddle)
	nk.NkStyleSetFont(ctx, avenirFont.Handle())

	if nk.NkButtonLabel(ctx, "Interface") == 1 {
		currentWindow = 3
	}
	if nk.NkButtonLabel(ctx, "Processing") == 1 {
		currentWindow = 4
	}

	nk.NkEnd(ctx)

	//   --------------------------------------------------
	//   Bottom panel
	//   --------------------------------------------------

	bounds = nk.NkRect(panelSize*winWidth+sideBarSize*winWidth, winHeight-200, winWidth-(2*panelSize*winWidth)-(sideBarSize*winWidth), 200)
	nk.NkBegin(ctx, "Files", bounds, nk.WindowNoScrollbar)

	nk.NkLayoutRowDynamic(ctx, 80, 1)
	nk.NkLabel(ctx, "Reee", nk.TextAlignCentered|nk.TextAlignMiddle)

	nk.NkEnd(ctx)

	nk.NkPlatformRender(nk.AntiAliasingOn, maxVertexBuffer, maxElementBuffer)
}

type Option uint8

const (
	Easy Option = 0
	Hard Option = 1
)

type State struct {
	bgColor nk.Color
	prop    int32
	opt     Option
}

func onError(code int32, msg string) {
	log.Printf("[glfw ERR]: error %d: %s", code, msg)
}

func b(v int32) bool {
	return v == 1
}

func flag(v bool) int32 {
	if v {
		return 1
	}
	return 0
}
