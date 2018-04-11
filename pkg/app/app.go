package app

import (
	g "gojira2d/pkg/graphics"
	"gojira2d/pkg/utils"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	OpenGLMajorVersion = 4
	OpenGLMinorVersion = 1
)

type App struct {
	Window     *glfw.Window
	Context    *g.Context
	FpsCounter *utils.FPSCounter
}

func InitApp(windowWidth int, windowHeight int, windowCentered bool, windowTitle string) *App {
	runtime.LockOSThread()
	app := &App{}
	app.Window = initWindow(windowWidth, windowHeight, windowTitle)
	app.Context = &g.Context{}
	app.Context.SetOrtho2DProjection(windowWidth, windowHeight, 1, windowCentered)
	app.FpsCounter = &utils.FPSCounter{}
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
	gl.ClearColor(0, 0, 0, 1.0)

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	return window
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

		a.Context.Clear()
		render()
		a.Context.RenderDrawableList()
		a.Context.EraseDrawableList()

		glfw.PollEvents()
		a.Window.SwapBuffers()
	}
}
