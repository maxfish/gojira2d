package app

import (
	"fmt"
	"gojira2d/pkg/graphics"
	g "gojira2d/pkg/graphics"
	"gojira2d/pkg/ui"
	"gojira2d/pkg/utils"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/go-gl/mathgl/mgl32"
)

const (
	OpenGLMajorVersion = 4
	OpenGLMinorVersion = 1
)

type App struct {
	Window         *glfw.Window
	Context        *g.Context
	UIContext      *g.Context
	FpsCounter     *utils.FPSCounter
	FpsCounterText *ui.Text

	clearColor g.Color
}

func InitApp(windowWidth int, windowHeight int, windowCentered bool, windowTitle string) *App {
	runtime.LockOSThread()
	app := &App{}
	app.Window = initWindow(windowWidth, windowHeight, windowTitle)
	app.Context = &g.Context{}
	app.Context.SetOrtho2DProjection(windowWidth, windowHeight, 1, windowCentered)
	app.UIContext = &g.Context{}
	app.UIContext.SetOrtho2DProjection(windowWidth, windowHeight, 1, false)
	app.FpsCounter = &utils.FPSCounter{}

	font := ui.NewFontFromFiles(
		"mono",
		"examples/assets/fonts/roboto-mono-regular.fnt",
		"examples/assets/fonts/roboto-mono-regular.png",
	)
	app.FpsCounterText = ui.NewText(
		"0",
		font,
		mgl32.Vec3{10, 10, -1},
		mgl32.Vec2{25, 25},
		graphics.Color{1, 0, 0, 1},
		mgl32.Vec4{0, 0, 0, -.17},
	)
	return app
}

func TerminateApp() {
	glfw.Terminate()
}

func initWindow(width, height int, title string) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, OpenGLMajorVersion)
	glfw.WindowHint(glfw.ContextVersionMinor, OpenGLMinorVersion)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// OpenGL
	if err := gl.Init(); err != nil {
		panic(err)
	}
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthMask(true)
	gl.DepthFunc(gl.LEQUAL)
	gl.DepthRange(0.0, 1.0)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	return window
}

// Clear clears the screen using App.clearColor
func (a *App) Clear() {
	gl.ClearColor(
		a.clearColor[0], a.clearColor[1], a.clearColor[2], a.clearColor[3])
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

// SetClearColor changes OpenGL background clear color
func (a *App) SetClearColor(color g.Color) {
	a.clearColor = color
}

func (a *App) MainLoop(
	update func(float64),
	render func(),
) {
	var newTime, oldTime float64
	for !a.Window.ShouldClose() {
		newTime = glfw.GetTime()
		deltaTime := newTime - oldTime
		update(deltaTime)
		a.FpsCounter.Update(deltaTime, 1)
		oldTime = newTime

		a.Clear()
		render()

		a.FpsCounterText.SetText(fmt.Sprintf("%v", a.FpsCounter.FPS()))
		a.UIContext.EnqueueForDrawing(a.FpsCounterText)
		a.UIContext.RenderDrawableList()
		a.UIContext.EraseDrawableList()

		a.Context.RenderDrawableList()
		a.Context.EraseDrawableList()

		glfw.PollEvents()
		a.Window.SwapBuffers()
	}
}
