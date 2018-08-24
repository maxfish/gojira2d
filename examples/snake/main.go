package main

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/maxfish/gojira2d/pkg/app"
	g "github.com/maxfish/gojira2d/pkg/graphics"
	"github.com/maxfish/gojira2d/pkg/input"
)

// Cell ...
type Cell struct {
	X, Y   int
	IsFood bool
}

// SnakeGame ...
type SnakeGame struct {
	directionX, directionY int
	cellSize, gridSize     int
	cellRadius             float32
	moveDelay, timePassed  float64
	worldBg, worldFg       *g.Primitive2D
	gridTopLeft            mgl32.Vec3
	snake                  []Cell
	food                   []Cell
}

// NewSnakeGame ...
func NewSnakeGame(cellSize, gridSize, width, height int) *SnakeGame {
	sg := &SnakeGame{
		cellSize:   cellSize,
		gridSize:   gridSize,
		directionX: 1,
		directionY: 0,
		moveDelay:  0.2,
		timePassed: 0.0,
		cellRadius: float32(cellSize) * math.Sqrt2 / 2,
		gridTopLeft: mgl32.Vec3{
			float32(width/2 - cellSize*gridSize/2),
			float32(height/2 - cellSize*gridSize/2),
			0,
		},
		snake: []Cell{Cell{0, 0, false}, Cell{0, 1, false}, Cell{0, 2, false}},
		food:  []Cell{Cell{3, 3, true}},
	}

	sg.worldBg = g.NewRegularPolygonPrimitive(
		mgl32.Vec3{float32(width / 2), float32(height / 2), 0},
		sg.cellRadius*float32(sg.gridSize)+2, 4, true)
	sg.worldBg.SetColor(g.Color{1, 1, 1, 1})
	sg.worldBg.SetAngle(45 * math.Pi / 180)

	sg.worldFg = g.NewRegularPolygonPrimitive(
		mgl32.Vec3{float32(width / 2), float32(height / 2), 0},
		sg.cellRadius*float32(sg.gridSize)+1, 4, true)
	sg.worldFg.SetColor(g.Color{0, 0, 0, 1})
	sg.worldFg.SetAngle(45 * math.Pi / 180)

	return sg
}

// GridToPos ...
func (sg *SnakeGame) GridToPos(gx, gy int) mgl32.Vec3 {
	return mgl32.Vec3{
		sg.gridTopLeft[0] + float32(gx*sg.cellSize+sg.cellSize/2),
		sg.gridTopLeft[1] + float32(gy*sg.cellSize+sg.cellSize/2),
		0,
	}
}

// DrawCell ...
func (sg *SnakeGame) DrawCell(cell Cell, ctx *g.Context) {
	p := g.NewRegularPolygonPrimitive(
		sg.GridToPos(cell.X, cell.Y), sg.cellRadius-1, 4, true)
	p.SetColor(g.Color{1, 1, 1, 1})
	p.SetAngle(45 * math.Pi / 180)
	p.Draw(app.Context)
}

// Draw ...
func (sg *SnakeGame) Draw(ctx *g.Context) {
	sg.worldBg.Draw(app.Context)
	sg.worldFg.Draw(app.Context)
	for _, cell := range sg.snake {
		sg.DrawCell(cell, app.Context)
	}
	for _, cell := range sg.food {
		sg.DrawCell(cell, app.Context)
	}
}

// UpdateDirection ...
func (sg *SnakeGame) UpdateDirection(kbd *input.KeyboardController) {
	kbd.Update()
	if kbd.ButtonPressed(input.ButtonDirPadUp) {
		if sg.directionY != 1 {
			sg.directionX = 0
			sg.directionY = -1
		}
	}
	if kbd.ButtonPressed(input.ButtonDirPadDown) {
		if sg.directionY != -1 {
			sg.directionX = 0
			sg.directionY = 1
		}
	}
	if kbd.ButtonPressed(input.ButtonDirPadLeft) {
		if sg.directionX != 1 {
			sg.directionX = -1
			sg.directionY = 0
		}
	}
	if kbd.ButtonPressed(input.ButtonDirPadRight) {
		if sg.directionX != -1 {
			sg.directionX = 1
			sg.directionY = 0
		}
	}
}

// Update ...
func (sg *SnakeGame) Update(kbd *input.KeyboardController, speed float64) {
	sg.UpdateDirection(kbd)

	sg.timePassed += speed
	if sg.timePassed > sg.moveDelay {
		var oldPos Cell
		for idx, pos := range sg.snake {
			if idx == 0 {
				oldPos = pos
				sg.snake[idx] = Cell{
					oldPos.X + sg.directionX,
					oldPos.Y + sg.directionY,
					false,
				}
			} else {
				tmpPos := pos
				sg.snake[idx] = oldPos
				oldPos = tmpPos
			}
		}
		sg.timePassed = 0
	}
}

func main() {
	app.Init(640, 480, "Snake")
	defer app.Terminate()
	app.SetFPSCounterVisible(true)

	// Use a cartesian coordinates system
	app.Context.Camera2D.SetFlipVertical(true)
	app.Context.Camera2D.SetPosition(50, 50)

	kbd := &input.KeyboardController{}
	kbd.Open(-1)

	game := NewSnakeGame(8, 22, 640, 480)

	app.MainLoop(
		func(speed float64) { game.Update(kbd, speed) },
		func() { game.Draw(app.Context) },
	)
}
