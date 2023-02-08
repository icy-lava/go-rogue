package main

import (
	"errors"
	// "fmt"
	"math"

	"github.com/nsf/termbox-go"
)

const width = 50
const height = 20

var grid [width][height]TileType
var entities []Entity
var playerMoved bool = false

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

func heuristic(from Position, to Position) int {
	dx := float64(to.x - from.x)
	dy := float64(to.y - from.y)
	result := int(math.Abs(dx) + math.Abs(dy))
	return result
}

func distanceToNeighbor(neighbor Position) int {
	switch grid[neighbor.x][neighbor.y] {
	case Empty:
		return 1
	case Wall:
		return 1000
	default:
		return 1
	}
}

func getScoreOrLarge(posMap map[Position]int, pos Position) int {
	score, ok := posMap[pos]
	if !ok {
		score = 1000000
	}
	return score
}

func findPath(start Position, end Position) ([]Position, error) {
	openSet := []Position{start}
	
	cameFrom := make(map[Position]*Position)
	gScore := make(map[Position]int)
	gScore[start] = 0
	
	fScore := make(map[Position]int)
	fScore[start] = heuristic(start, end)
	
	for len(openSet) != 0 {
		currentI := 0
		current := openSet[currentI]
		for i := 1; i < len(openSet); i++ {
			node := openSet[i]
			currentScore := getScoreOrLarge(fScore, current)
			nodeScore := getScoreOrLarge(fScore, node)
			if nodeScore < currentScore {
				current = node
				currentI = i
			}
		}
		if current == end {
			totalPath := []Position{current}
			for {
				next := cameFrom[current]
				if next == nil {break}
				current = *next
				totalPath = append(totalPath, current)
			}
			for i, j := 0, len(totalPath)-1; i < j; i, j = i+1, j-1 {
				totalPath[i], totalPath[j] = totalPath[j], totalPath[i]
			}
			return totalPath, nil
		}
		openSet[currentI] = openSet[len(openSet)-1]
		openSet = openSet[:len(openSet)-1]
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {continue}
				
				pos := Position{
					x: current.x + dx,
					y: current.y + dy,
				}
				if pos.x < 0 || pos.y < 0 || pos.x >= width || pos.y >= height {
					continue
				}
				
				if grid[pos.x][pos.y] == Wall {continue}
				
				neighbor := pos
				tentativeGScore := getScoreOrLarge(gScore, current) + distanceToNeighbor(neighbor)
				neighborScore := getScoreOrLarge(gScore, neighbor)
				if tentativeGScore < neighborScore {
					from := current
					cameFrom[neighbor] = &from
					gScore[neighbor] = tentativeGScore
					fScore[neighbor] = tentativeGScore + heuristic(neighbor, end)
					var found bool
					for _, node := range openSet {
						if node == neighbor {
							found = true
							break
						}
					}
					if !found {
						openSet = append(openSet, neighbor)
					}
				}
			}
		}
	}
	return nil, errors.New("could not find path")
}

type Entity interface {
	GetSymbol() rune
	
	GetPosition() (int, int)
	SetPosition(x int, y int)
	
	GetTargetPosition() Position
	Move(dx int, dy int)
}

type Controlled interface {
	SetLastInput(dx int, dy int)
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

type Rat struct {
	Position
}

func (entity Rat) GetSymbol() rune {
	return 'r'
}

func (entity Rat) GetTargetPosition() Position {
	var nearestPlayer Entity
	for _, entity := range entities {
		switch entity.(type) {
		case *Player:
			// TODO: find nearest player
			nearestPlayer = entity
		}
	}
	if nearestPlayer != nil {
		x, y := nearestPlayer.GetPosition()
		return Position{x, y}
	}
	x, y := entity.GetPosition()
	return Position{x, y}
}

// TODO: remove code duplication
func (entity *Rat) Move(dx int, dy int) {
	if math.Abs(float64(dx)) > 1 || math.Abs(float64(dy)) > 1 {
		panic("moving further than 1 tile")
	}
	tile := grid[entity.x + dx][entity.y + dy]
	if tile != Wall {
		entity.x += dx
		entity.y += dy
	}
}

type Player struct {
	Position
	lastInputDX int
	lastInputDY int
}

func (entity Player) GetSymbol() rune {
	return '@'
}

func (entity Player) GetTargetPosition() Position {
	return Position{
		x: entity.x + entity.lastInputDX,
		y: entity.y + entity.lastInputDY,
	}
}

func (entity *Player) SetLastInput(dx int, dy int) {
	entity.lastInputDX = dx
	entity.lastInputDY = dy
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
				termbox.SetChar(x, y, '▓')
				// termbox.SetCell(x, y, '▓', termbox.ColorBlue, termbox.ColorDefault)
			}
		}
	}
	for _, entity := range entities {
		symbol := entity.GetSymbol()
		x, y := entity.GetPosition()
		termbox.SetChar(x, y, symbol)
	}
}

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
		Position: Position{
			x: 20,
			y: 10,
		},
		lastInputDX: 0,
		lastInputDY: 0,
	})
	
	entities = append(entities, &Rat{
		Position: Position{
			x: 30,
			y: 12,
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
					controlled, ok := entity.(Controlled)
					if ok {
						controlled.SetLastInput(dx, dy)
						// NOTE: might not be a player though
						playerMoved = true
					}
				}
			}
			if playerMoved {
				playerMoved = false
				
				targetMap := make(map[Entity]Position)
				for _, entity := range entities {
					targetMap[entity] = entity.GetTargetPosition()
				}
				
				for entity, tpos := range targetMap {
					ex, ey := entity.GetPosition()
					epos := Position{ex, ey}
					// fmt.Printf("%v %v\n", epos, tpos)
					if tpos != epos {
						path, err := findPath(epos, tpos)
						if err == nil {
							nextPos := path[1]
							entity.Move(nextPos.x - epos.x, nextPos.y - epos.y)
						}
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
