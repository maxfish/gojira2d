package input

import (
	"github.com/maxfish/gojira2d/pkg/graphics"
)

func UpdateCameraWithMouse(camera *graphics.Camera2D) {
	deltaX, deltaY := MouseDelta()
	if MouseButton(1) {
		tX := (float64(-deltaX)) / camera.Zoom()
		tY := (float64(-deltaY)) / camera.Zoom()
		camera.Translate(tX, tY)
		//fmt.Printf("X:%f  Y:%f\n", tX, tY)
	}
	_, wY := MouseScroll()
	camera.SetZoom(camera.Zoom() + wY/10)
}
