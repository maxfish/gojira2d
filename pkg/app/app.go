package app

import (
	"log"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v4.1-core/gl"
	"runtime"
	g "gojira2d/pkg/graphics"
)

const (
	OpenGLMajorVersion = 4
	OpenGLMinorVersion = 1
)

type App struct {
	Window  *glfw.Window
	Context *g.Context
}

func InitApp(windowWidth int, windowHeight int, windowCentered bool, windowTitle string) (*App) {
	runtime.LockOSThread()
	core := &App{}
	core.Window = initWindow(windowWidth, windowHeight, windowTitle)
	core.Context = &g.Context{}
	core.Context.SetOrtho2DProjection(windowWidth, windowHeight, 1, windowCentered)
	return core
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
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.ClearColor(0, 0, 0, 1.0)

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	return window
}

func (c *App) MainLoop(
	update func(float64),
	render func(),
) {
	var newTime, oldTime float64
	for !c.Window.ShouldClose() {
		newTime = glfw.GetTime()
		update(newTime - oldTime)
		oldTime = newTime

		c.Context.Clear()
		render()
		c.Context.RenderDrawableList()
		c.Context.EraseDrawableList()

		glfw.PollEvents()
		c.Window.SwapBuffers()
	}
}
