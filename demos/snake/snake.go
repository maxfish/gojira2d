package main

import (
	"math"
	"math/rand"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/maxfish/gojira2d/pkg/app"
	g "github.com/maxfish/gojira2d/pkg/graphics"
	"github.com/maxfish/gojira2d/pkg/input"
)

type cell struct {
	x, y   int
	isFood bool
}

type direction int

const (
	dirUp    direction = 1
	dirDown            = -1
	dirRight           = 2
	dirLeft            = -2
)

type snakeGame struct {
	direction             direction
	cellSize, gridSize    int
	cellRadius            float64
	moveDelay, timePassed float64
	gridTopLeft           mgl64.Vec3
	snake                 []cell
	food                  []bool
	initScreen, spawnFood float64
	speedupAfter          float64
}

func newSnakeGame(cellSize, gridSize, width, height int) *snakeGame {
	sg := &snakeGame{
		cellSize:   cellSize,
		gridSize:   gridSize,
		direction:  dirUp,
		moveDelay:  0.2,
		timePassed: 0,
		cellRadius: float64(cellSize/2) * math.Sqrt2,
		gridTopLeft: mgl64.Vec3{
			float64(width/2 - cellSize*gridSize/2),
			float64(height/2 - cellSize*gridSize/2),
			0,
		},
		snake:        []cell{cell{0, 2, false}, cell{0, 1, false}, cell{0, 0, false}},
		food:         make([]bool, gridSize*gridSize),
		initScreen:   3,
		spawnFood:    5,
		speedupAfter: 10,
	}

	return sg
}

func (sg *snakeGame) cellToWorld(gx, gy int) mgl64.Vec3 {
	return mgl64.Vec3{
		sg.gridTopLeft[0] + float64(gx*sg.cellSize+sg.cellSize/2),
		sg.gridTopLeft[1] + float64(gy*sg.cellSize+sg.cellSize/2),
		0,
	}
}

func (sg *snakeGame) drawCell(cell cell, ctx *g.Context) {
	p := g.NewRegularPolygonPrimitive(
		sg.cellToWorld(cell.x, cell.y), sg.cellRadius-1, 4, true)
	p.SetColor(g.Color{1, 1, 1, 1})
	p.SetAngle(45 * math.Pi / 180)
	p.Draw(app.Context)
}

func (sg *snakeGame) isValidPos(x, y int) bool {
	return x < sg.gridSize && y < sg.gridSize && x >= 0 && y >= 0
}

func (sg *snakeGame) foodIdx(x, y int) int {
	if !sg.isValidPos(x, y) {
		panic("")
	}
	return y*sg.gridSize + x
}

func (sg *snakeGame) isFood(x, y int) bool {
	return sg.food[sg.foodIdx(x, y)]
}

func (sg *snakeGame) isSnake(x, y int) bool {
	for _, cell := range sg.snake {
		if cell.x == x && cell.y == y {
			return true
		}
	}
	return false
}

func (sg *snakeGame) draw() {
	for y := 0; y < sg.gridSize; y++ {
		for x := 0; x < sg.gridSize; x++ {
			if sg.isFood(x, y) {
				sg.drawCell(cell{x, y, true}, app.Context)
			}
		}
	}

	if sg.initScreen > 0 {
		return
	}

	for _, cell := range sg.snake {
		sg.drawCell(cell, app.Context)
	}
}

func (sg *snakeGame) snakeLayout() direction {
	head := sg.snake[0]
	neck := sg.snake[1]
	if head.x == neck.x {
		return direction(head.y - neck.y)
	}
	return direction((head.x - neck.x) * 2)
}

func (sg *snakeGame) updateDirection(kbd *input.KeyboardController) {
	kbd.Update()
	var newDir direction
	if kbd.ButtonPressed(input.ButtonDirPadDown) {
		newDir = dirDown
	} else if kbd.ButtonPressed(input.ButtonDirPadUp) {
		newDir = dirUp
	} else if kbd.ButtonPressed(input.ButtonDirPadLeft) {
		newDir = dirLeft
	} else if kbd.ButtonPressed(input.ButtonDirPadRight) {
		newDir = dirRight
	}
	if newDir != 0 && newDir != -sg.snakeLayout() {
		sg.direction = newDir
	}
}

func (sg *snakeGame) updatePosition(speed float64) {
	gameEnd := false
	sg.timePassed += speed
	if sg.timePassed > sg.moveDelay {
		var oldPos, headPos cell
		for idx, pos := range sg.snake {
			if idx == 0 {
				oldPos = pos
				var moveDirX, moveDirY int = 0, 0
				if math.Abs(float64(sg.direction)) > 1.5 {
					moveDirX = int(sg.direction) / 2
				} else {
					moveDirY = int(sg.direction)
				}
				headPos = cell{oldPos.x + moveDirX, oldPos.y + moveDirY, false}

				if !sg.isValidPos(headPos.x, headPos.y) || sg.isSnake(headPos.x, headPos.y) {
					gameEnd = true
					break
				}

				sg.snake[0] = headPos

			} else {
				tmpPos := pos
				sg.snake[idx] = oldPos
				oldPos = tmpPos
			}
		}
		if gameEnd {
			sg.initScreen = 3
		} else {
			head := sg.snake[0]
			if sg.isFood(head.x, head.y) {
				sg.snake = append(sg.snake, oldPos)
				sg.food[sg.foodIdx(head.x, head.y)] = false
			}
		}
		sg.timePassed = 0
	}
}

func (sg *snakeGame) updateFood() {
	idxArr := make([]int, 0)
	for y := 0; y < sg.gridSize; y++ {
		for x := 0; x < sg.gridSize; x++ {
			foodIdx := sg.foodIdx(x, y)
			if sg.food[foodIdx] {
				return
			}
			if sg.isSnake(x, y) {
				continue
			}
			idxArr = append(idxArr, foodIdx)
		}
	}

	if len(idxArr) == 0 {
		return
	}

	idxIdx := rand.Intn(len(idxArr))
	sg.food[idxArr[idxIdx]] = true
}

func (sg *snakeGame) updateSpeed(speed float64) {
	sg.speedupAfter -= speed
	if sg.speedupAfter < 0 {
		sg.speedupAfter = 10
		sg.moveDelay *= 0.7
	}
}

func (sg *snakeGame) update(kbd *input.KeyboardController, speed float64) {
	if sg.initScreen > 0 {
		sg.initScreen -= speed
		if sg.initScreen < 0 {
			sg.initScreen = 0
		}
		for i := 0; i < sg.gridSize*sg.gridSize; i++ {
			sg.food[i] = int(sg.initScreen) > i
		}
		return
	}

	sg.updateDirection(kbd)
	sg.updatePosition(speed)
	sg.updateFood()
	sg.updateSpeed(speed)

	if sg.initScreen > 0 {
		sg.snake = []cell{cell{0, 2, false}, cell{0, 1, false}, cell{0, 0, false}}
		sg.food = make([]bool, sg.gridSize*sg.gridSize)
		sg.direction = dirUp
		sg.moveDelay = 0.2
		sg.speedupAfter = 10
	}
}
