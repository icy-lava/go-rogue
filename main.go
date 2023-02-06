package main

import (
	"math"

	"github.com/nsf/termbox-go"
)

const width = 50
const height = 20

var grid [width][height]TileType
var entities []Entity

func assert1(err error) {
	if err != nil {
		panic(err)
	}
}

func assert2[V any](value V, err error) V {
	if err != nil {
		panic(err)
	}
	return value
}

type Entity interface {
	GetSymbol() rune
	
	GetPosition() (int, int)
	SetPosition(x int, y int)
	
	Move(dx int, dy int)
	TakesInput() bool
}

type Position struct {
	x int
	y int
}

func (pos Position) GetPosition() (int, int) {
	return pos.x, pos.y
}
func (pos *Position) SetPosition(x int, y int) {
	pos.x, pos.y = x, y
}

type Player struct {
	Position
}

func (entity *Player) GetSymbol() rune {
	return '@'
}
func (entity *Player) Move(dx int, dy int) {
	if math.Abs(float64(dx)) > 1 || math.Abs(float64(dy)) > 1 {
		panic("moving further than 1 tile")
	}
	tile := grid[entity.x + dx][entity.y + dy]
	if tile != Wall {
		entity.x += dx
		entity.y += dy
	}
}
func (entity *Player) TakesInput() bool {
	return true
}

type TileType int
const (
	Empty TileType = iota
	Wall
)

func redrawGrid() {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			switch grid[x][y] {
			case Empty:
				termbox.SetChar(x, y, '.')
			case Wall:
				termbox.SetChar(x, y, '#')
			}
		}
	}
	for _, entity := range entities {
		symbol := entity.GetSymbol()
		x, y := entity.GetPosition()
		termbox.SetChar(x, y, symbol)
	}
}

func noop(value any) {}

func main() {
	assert1(termbox.Init())
	defer termbox.Close()

	// termbox.Sync()
	
	for i := 10; i <= 40; i++ {
		grid[i][5] = Wall
		grid[i][15] = Wall
	}
	for i := 5; i <= 15; i++ {
		grid[10][i] = Wall
		grid[40][i] = Wall
	}
	
	entities = append(entities, &Player{
		Position{
			x: 20,
			y: 10,
		},
	})

	termbox.SetInputMode(termbox.InputEsc)
	assert1(termbox.Clear(termbox.ColorDefault, termbox.ColorDefault))
	redrawGrid()
	termbox.Flush()
	
game:
	for {
		switch event := termbox.PollEvent(); event.Type {
		case termbox.EventKey:
			if event.Key == termbox.KeyCtrlC {
				break game
			}
			
			// Player movement
			dx := 0
			dy := 0
			if event.Key == termbox.KeyArrowRight {
				dx++
			}
			if event.Key == termbox.KeyArrowLeft {
				dx--
			}
			if event.Key == termbox.KeyArrowDown {
				dy++
			}
			if event.Key == termbox.KeyArrowUp {
				dy--
			}
			if dx != 0 || dy != 0 {
				for _, entity := range entities {
					if entity.TakesInput() {
						entity.Move(dx, dy)
					}
				}
				redrawGrid()
				termbox.Flush()
			}
		case termbox.EventError:
			panic(event.Err)
		}
	}
}
