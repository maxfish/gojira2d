package app

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/maxfish/gojira2d/pkg/graphics"
	g "github.com/maxfish/gojira2d/pkg/graphics"
	"github.com/maxfish/gojira2d/pkg/ui"
	"github.com/maxfish/gojira2d/pkg/utils"
)

const (
	OpenGLMajorVersion = 4
	OpenGLMinorVersion = 1
)

var (
	window         *glfw.Window
	Context        *g.Context
	UIContext      *g.Context
	FpsCounter     *utils.FPSCounter
	FpsCounterText *ui.Text

	clearColor g.Color
)

func init() {
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		panic(err)
	}
}

func Init(windowWidth int, windowHeight int, windowCentered bool, windowTitle string) {
	window = initWindow(windowWidth, windowHeight, windowTitle)
	Context = &g.Context{}
	Context.Camera2D = g.NewCamera2D(windowWidth, windowHeight, 1, windowCentered)
	UIContext = &g.Context{}
	UIContext.Camera2D = g.NewCamera2D(windowWidth, windowHeight, 1, false)
	FpsCounter = &utils.FPSCounter{}

	font := ui.NewFontFromFiles(
		"mono",
		"examples/assets/fonts/roboto-mono-regular.fnt",
		"examples/assets/fonts/roboto-mono-regular.png",
	)
	FpsCounterText = ui.NewText(
		"0",
		font,
		mgl32.Vec3{float32(windowWidth - 30), 10, -1},
		mgl32.Vec2{25, 25},
		graphics.Color{1, 0, 0, 1},
		mgl32.Vec4{0, 0, 0, -.17},
	)
}

func Terminate() {
	glfw.Terminate()
}

func initWindow(width, height int, title string) *glfw.Window {
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

func GetWindow() *glfw.Window {
	return window
}

// Clear clears the screen using App.clearColor
func Clear() {
	gl.ClearColor(
		clearColor[0], clearColor[1], clearColor[2], clearColor[3])
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

// SetClearColor changes OpenGL background clear color
func SetClearColor(color g.Color) {
	clearColor = color
}

func MainLoop(
	update func(float64),
	render func(),
) {
	var newTime, oldTime, deltaTime float64
	for !window.ShouldClose() {
		newTime = glfw.GetTime()
		deltaTime = newTime - oldTime
		oldTime = newTime

		update(deltaTime)
		FpsCounter.Update(deltaTime, 1)
		FpsCounterText.SetText(fmt.Sprintf("%v", FpsCounter.FPS()))

		Clear()
		render()

		FpsCounterText.Draw(UIContext)

		glfw.PollEvents()
		window.SwapBuffers()
	}
}
