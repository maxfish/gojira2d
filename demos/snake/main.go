package main

import (
	"github.com/maxfish/gojira2d/pkg/app"
	"github.com/maxfish/gojira2d/pkg/input"
)

func main() {
	worldSize := 800
	cellSize := 20
	gridSize := worldSize / cellSize

	app.Init(worldSize, worldSize, "Snake")
	defer app.Terminate()

	app.SetFPSCounterVisible(true)
	app.Context.Camera2D.SetFlipVertical(true)
	app.Context.Camera2D.SetPosition(float64(worldSize/2), float64(worldSize/2))
	app.Context.Camera2D.SetVisibleArea(0, 0, float32(worldSize), float32(worldSize))

	kbd := &input.KeyboardController{}
	kbd.Open(-1)

	game := newSnakeGame(cellSize, gridSize, worldSize, worldSize)
	update := func(deltaTime float64) { game.update(kbd, deltaTime) }
	app.MainLoop(update, game.draw)
}
