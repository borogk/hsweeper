package game

import (
	"math/rand"
)

type (
	Status int

	RevealResult int

	Point struct {
		x int
		y int
	}

	Game struct {
		status              Status
		cells               []Cell
		width               int
		height              int
		mines               int
		hearts              int
		livesLeft           int
		heartsLeft          int
		revealedCounter     int
		flaggedCounter      int
		heartSpawnCounter   int
		heartSpawnThreshold int
	}
)

const (
	StatusReady Status = iota
	StatusStarted
	StatusLost
	StatusWon
)

const (
	RevealResultOk RevealResult = iota
	RevealResultBlast
	RevealResultBlocked
)

var AdjacentOffsets = []Point{
	{x: -1, y: -1},
	{x: 0, y: -1},
	{x: 1, y: -1},
	{x: 1, y: 0},
	{x: 1, y: 1},
	{x: 0, y: 1},
	{x: -1, y: 1},
	{x: -1, y: 0},
}

func NewGame(width, height, mines, hearts, livesLeft int) *Game {
	return &Game{
		status:     StatusReady,
		cells:      make([]Cell, width*height),
		width:      width,
		height:     height,
		mines:      mines,
		hearts:     hearts,
		livesLeft:  livesLeft,
		heartsLeft: hearts,
	}
}

func (g *Game) Status() Status {
	return g.status
}

func (g *Game) Width() int {
	return g.width
}

func (g *Game) Height() int {
	return g.height
}

func (g *Game) IsOutOfBounds(x, y int) bool {
	return x < 0 || x >= g.width || y < 0 || y >= g.height
}

func (g *Game) Cell(x, y int) *Cell {
	if !g.IsOutOfBounds(x, y) {
		i := x + y*g.width
		return &g.cells[i]
	}

	return &Cell{isRevealed: true}
}

func (g *Game) IsFinished() bool {
	return g.status == StatusLost || g.status == StatusWon
}

func (g *Game) MinesRemaining() int {
	return g.mines - g.flaggedCounter
}

func (g *Game) LivesRemaining() int {
	return g.livesLeft
}

func (g *Game) ToggleFlag(x, y int) {
	if g.IsFinished() {
		return
	}

	cell := g.Cell(x, y)
	if !cell.isRevealed {
		cell.isFlagged = !cell.isFlagged
		cell.isQuestioned = false
		if cell.isFlagged {
			g.flaggedCounter++
		} else {
			g.flaggedCounter--
		}
	}
}

func (g *Game) ToggleQuestion(x, y int) {
	if g.IsFinished() {
		return
	}

	cell := g.Cell(x, y)
	if !cell.isRevealed {
		cell.isQuestioned = !cell.isQuestioned
		if cell.isFlagged {
			cell.isFlagged = false
			g.flaggedCounter--
		}
	}
}

func (g *Game) ClearFlagAndQuestion(x, y int) {
	if g.IsFinished() {
		return
	}

	cell := g.Cell(x, y)
	if cell.isFlagged {
		cell.isFlagged = false
		g.flaggedCounter--
	}
	cell.isQuestioned = false
}

func (g *Game) Pickup(x, y int) {
	if g.IsFinished() {
		return
	}

	cell := g.Cell(x, y)
	if cell.isHeart {
		cell.isHeart = false
		g.livesLeft++
	}
}

func (g *Game) Reveal(x, y int) RevealResult {
	cell := g.Cell(x, y)
	if cell.isRevealed || cell.isFlagged || cell.isQuestioned {
		return RevealResultBlocked
	}

	if g.status == StatusReady {
		mineLocations := g.randomMineLocations(x, y)
		g.plantMines(mineLocations)
		g.calculateHeartSpawnThreshold()
		g.status = StatusStarted
	}

	cell.isRevealed = true
	g.revealedCounter++

	result := RevealResultOk
	if cell.isMine {
		result = RevealResultBlast
		g.livesLeft--
		if g.livesLeft <= 0 {
			g.status = StatusLost
		} else {
			cell.isMine = false
			g.mines--
			for _, offset := range AdjacentOffsets {
				g.Cell(x+offset.x, y+offset.y).adjacentMines--
			}
		}
	} else if len(g.cells)-g.revealedCounter == g.mines {
		g.status = StatusWon
	}

	if cell.adjacentMines == 0 {
		g.heartSpawnCounter++
		if g.heartsLeft > 0 && g.heartSpawnCounter%g.heartSpawnThreshold == 0 {
			cell.isHeart = true
			g.heartsLeft--
		}

		for _, offset := range AdjacentOffsets {
			subResult := g.Reveal(x+offset.x, y+offset.y)
			if subResult == RevealResultBlast {
				result = RevealResultBlast
			}
		}
	}

	return result
}

func (g *Game) RevealAdjacent(x, y int) RevealResult {
	cell := g.Cell(x, y)
	if !cell.isRevealed || cell.adjacentMines == 0 {
		return RevealResultBlocked
	}

	adjacentFlags := 0
	for _, offset := range AdjacentOffsets {
		if g.Cell(x+offset.x, y+offset.y).isFlagged {
			adjacentFlags++
		}
	}

	result := RevealResultOk
	if adjacentFlags == cell.adjacentMines {
		for _, offset := range AdjacentOffsets {
			subResult := g.Reveal(x+offset.x, y+offset.y)
			if subResult == RevealResultBlast {
				result = RevealResultBlast
			}
		}
	}

	return result
}

func (g *Game) randomMineLocations(aroundX, aroundY int) []Point {
	locations := make([]Point, 0, g.mines)
	for _, i := range rand.Perm(len(g.cells)) {
		if len(locations) == g.mines {
			break
		}

		x := i % g.width
		y := i / g.width
		if x < aroundX-1 || x > aroundX+1 || y < aroundY-1 || y > aroundY+1 {
			locations = append(locations, Point{x, y})
		}
	}

	return locations
}

func (g *Game) plantMines(mineLocations []Point) {
	for _, location := range mineLocations {
		g.Cell(location.x, location.y).isMine = true
	}

	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			cell := g.Cell(x, y)
			cell.adjacentMines = 0
			for _, offset := range AdjacentOffsets {
				if g.Cell(x+offset.x, y+offset.y).isMine {
					cell.adjacentMines++
				}
			}
		}
	}
}

func (g *Game) calculateHeartSpawnThreshold() {
	empty := 0
	for _, cell := range g.cells {
		if cell.adjacentMines == 0 && !cell.isMine {
			empty++
		}
	}

	g.heartSpawnThreshold = empty / (g.hearts + 1)
	if g.heartSpawnThreshold == 0 {
		g.heartSpawnThreshold = 1
	}
}
