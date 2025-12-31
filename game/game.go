package game

import (
	"math/rand"
	"sync"
)

// Game encapsulates a game of hsweeper with its entire logic.
type Game struct {
	status              Status
	cells               []Cell
	width               int
	height              int
	minesToPlant        int
	heartsToPlant       int
	livesLeft           int
	minesLeft           int
	heartsLeft          int
	unrevealedCounter   int
	flaggedCounter      int
	heartSpawnCounter   int
	heartSpawnThreshold int
	sync.Mutex
}

// NewGame creates a new game with desired parameters.
func NewGame(width, height, minesToPlant, heartsToPlant, livesLeft int) *Game {
	// Validate and auto-correct parameters
	status := StatusReady
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	if minesToPlant < 0 {
		minesToPlant = 0
	}
	if heartsToPlant < 0 {
		heartsToPlant = 0
	}

	// Allow zero initial lives, but auto-lose for consistency
	if livesLeft < 1 {
		livesLeft = 0
		status = StatusLost
	}

	return &Game{
		status:            status,
		cells:             make([]Cell, width*height),
		width:             width,
		height:            height,
		minesToPlant:      minesToPlant,
		heartsToPlant:     heartsToPlant,
		livesLeft:         livesLeft,
		heartsLeft:        heartsToPlant,
		unrevealedCounter: width * height,
	}
}

// RestoreGame creates a game and restores it to the state as told by provided snapshot.
// Prioritizes creating a playable game with consistent state over being 100% faithful to the snapshot parameters.
func RestoreGame(snapshot *Snapshot) *Game {
	game := NewGame(snapshot.Width, snapshot.Height, snapshot.MinesToPlant, snapshot.HeartsToPlant, snapshot.LivesLeft)

	// Ignore the rest of snapshot parameters if the game wasn't supposed to start yet
	if snapshot.Status == StatusReady {
		return game
	}

	// Set the game status but don't override lost, which can happen on a snapshot with 0 lives left
	if game.status != StatusLost {
		game.status = snapshot.Status
	}

	// Plant mines
	game.plantMines(snapshot.MineLocations)

	// Reveal cells one by one with consistency checks
	for _, i := range snapshot.RevealedLocations {
		cell := &game.cells[i]
		if !cell.isRevealed && !cell.isMine {
			cell.isRevealed = true
			game.unrevealedCounter--

			// Keep track of heart spawn stats
			if cell.adjacentMines == 0 {
				game.heartSpawnCounter++
			}
		}
	}

	// Set hearts left and plant uncollected hearts with consistency checks
	game.heartsLeft = snapshot.HeartsLeft
	for _, i := range snapshot.UncollectedHeartLocations {
		cell := &game.cells[i]
		if !cell.isHeart && cell.isRevealed && cell.adjacentMines == 0 {
			cell.isHeart = true
		}
	}

	// Put flags with consistency checks
	for _, i := range snapshot.FlaggedLocations {
		cell := &game.cells[i]
		if !cell.isFlagged && !cell.isRevealed {
			cell.isFlagged = true
			game.flaggedCounter++
		}
	}

	// Put questions with consistency checks
	for _, i := range snapshot.QuestionedLocations {
		cell := &game.cells[i]
		if !cell.isFlagged && !cell.isRevealed {
			cell.isQuestioned = true
		}
	}

	return game
}

// Save captures current game state into a snapshot.
func (g *Game) Save() *Snapshot {
	g.Lock()
	defer g.Unlock()

	return &Snapshot{
		Status:                    g.status,
		Width:                     g.width,
		Height:                    g.height,
		MinesToPlant:              g.minesToPlant,
		HeartsToPlant:             g.heartsToPlant,
		LivesLeft:                 g.livesLeft,
		HeartsLeft:                g.heartsLeft,
		MineLocations:             g.collectLocationsForSnapshot(isCellMine),
		RevealedLocations:         g.collectLocationsForSnapshot(isCellRevealed),
		UncollectedHeartLocations: g.collectLocationsForSnapshot(isCellHeart),
		FlaggedLocations:          g.collectLocationsForSnapshot(isCellFlagged),
		QuestionedLocations:       g.collectLocationsForSnapshot(isCellQuestioned),
	}
}

// Status returns game status
func (g *Game) Status() Status {
	return g.status
}

// Width returns game width
func (g *Game) Width() int {
	return g.width
}

// Height returns game height
func (g *Game) Height() int {
	return g.height
}

// IsOutOfBounds checks if provided coordinates are out of bounds.
func (g *Game) IsOutOfBounds(x, y int) bool {
	return x < 0 || x >= g.width || y < 0 || y >= g.height
}

// Cell returns pointer to a cell.
func (g *Game) Cell(x, y int) *Cell {
	if !g.IsOutOfBounds(x, y) {
		i := x + y*g.width
		return &g.cells[i]
	}

	// Having out-of-bounds representation prevents crashes as well as stops recursive reveals from propagating.
	return &Cell{isRevealed: true}
}

// IsFinished indicates if the game is one of the finished states.
func (g *Game) IsFinished() bool {
	return g.status == StatusLost || g.status == StatusWon
}

// MinesRemaining indicates how many mines are left, taking flags into account.
func (g *Game) MinesRemaining() int {
	if g.status == StatusReady {
		return g.minesToPlant
	}

	return g.minesLeft - g.flaggedCounter
}

// LivesRemaining indicates how many lives are left.
func (g *Game) LivesRemaining() int {
	return g.livesLeft
}

// ToggleFlag toggles flagged state of a cell, removes question.
func (g *Game) ToggleFlag(x, y int) {
	g.Lock()
	defer g.Unlock()

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

// ToggleQuestion toggles questioned state of a cell, removes flag.
func (g *Game) ToggleQuestion(x, y int) {
	g.Lock()
	defer g.Unlock()

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

// ClearFlagAndQuestion clears all markings on a cell.
func (g *Game) ClearFlagAndQuestion(x, y int) {
	g.Lock()
	defer g.Unlock()

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

// Pickup collects a heart if there is one.
func (g *Game) Pickup(x, y int) {
	g.Lock()
	defer g.Unlock()

	if g.IsFinished() {
		return
	}

	cell := g.Cell(x, y)
	if cell.isHeart {
		cell.isHeart = false
		g.livesLeft++
	}
}

// Reveal reveals a cell, advancing the game forward.
// Recursively propagates over cells with 0 adjacent mines.
func (g *Game) Reveal(x, y int) RevealResult {
	g.Lock()
	defer g.Unlock()
	return g.revealInner(x, y)
}

// AdvancedReveal reveals adjacent cells after exact amount of them were flagged.
func (g *Game) AdvancedReveal(x, y int) RevealResult {
	g.Lock()
	defer g.Unlock()

	cell := g.Cell(x, y)

	// Only allow revealed numbered cells
	if !cell.isRevealed || cell.adjacentMines == 0 {
		return RevealResultBlocked
	}

	// Proceed only if adjacent flags match the cell number
	adjacentFlags := 0
	for _, offset := range adjacentOffsets {
		if g.Cell(x+offset.x, y+offset.y).isFlagged {
			adjacentFlags++
		}
	}
	if adjacentFlags != cell.adjacentMines {
		return RevealResultBlocked
	}

	// Result types are ordered Blocked-Revealed-Blast, so treat the maximum as the combined result
	// If any were revealed - combined result would be at least revealed, if any blasted - blast
	result := RevealResultBlocked
	for _, offset := range adjacentOffsets {
		subResult := g.revealInner(x+offset.x, y+offset.y)
		if subResult > result {
			result = subResult
		}
	}

	// Blast means the adjacent flags were incorrect, remove them for safety
	if result == RevealResultBlast {
		for _, offset := range adjacentOffsets {
			adjacentCell := g.Cell(x+offset.x, y+offset.y)
			if adjacentCell.isFlagged {
				adjacentCell.isFlagged = false
				g.flaggedCounter--
			}
		}
	}

	return result
}

// Inner implementation of Reveal, extracted to avoid locking twice on recursion.
func (g *Game) revealInner(x, y int) RevealResult {
	if g.IsFinished() {
		return RevealResultBlocked
	}

	cell := g.Cell(x, y)

	// Only allow unrevealed and unmarked cells
	if cell.isRevealed || cell.isFlagged || cell.isQuestioned {
		return RevealResultBlocked
	}

	// First reveal triggers game initialization
	if g.status == StatusReady {
		g.plantMines(g.randomMineLocations(x, y))
		g.status = StatusStarted
	}

	// Mark as revealed
	result := RevealResultRevealed
	cell.isRevealed = true
	g.unrevealedCounter--

	// Check if we hit a mine
	if cell.isMine {
		result = RevealResultBlast
		g.livesLeft--

		if g.livesLeft > 0 {
			// Some lives left, remove the mine and adjust neighboring cell numbers
			cell.isMine = false
			g.minesLeft--
			for _, offset := range adjacentOffsets {
				adjacentCell := g.Cell(x+offset.x, y+offset.y)
				adjacentCell.adjacentMines--
				// After blast an adjacent cell might become eligible for propagation
				if adjacentCell.isRevealed && adjacentCell.adjacentMines == 0 {
					g.propagateReveal(x+offset.x, y+offset.y)
				}
			}
		} else {
			// No lives left, declare loss
			g.status = StatusLost
		}
	}

	// Propagate reveal
	if cell.adjacentMines == 0 {
		g.propagateReveal(x, y)
	}

	// We are still alive and only mines are left unrevealed, declare victory
	if g.status != StatusLost && g.unrevealedCounter == g.minesLeft {
		g.status = StatusWon
	}

	return result
}

// Reveals adjacent cells, center of propagation itself must be revealed and isolated.
// This function is called once as soon as both conditions are met.
func (g *Game) propagateReveal(x, y int) {
	cell := g.Cell(x, y)

	// Process heart spawning here, as revealed isolated cells are all candidates for having a pickup
	g.heartSpawnCounter++
	if g.heartsLeft > 0 && g.heartSpawnCounter%g.heartSpawnThreshold == 0 {
		cell.isHeart = true
		g.heartsLeft--
	}

	for _, offset := range adjacentOffsets {
		g.revealInner(x+offset.x, y+offset.y)
	}
}

// Collects list of indices of cells matching specified predicate.
func (g *Game) collectLocationsForSnapshot(predicate CellPredicate) []int {
	locations := make([]int, 0)
	for i, cell := range g.cells {
		if predicate(&cell) {
			locations = append(locations, i)
		}
	}

	if len(locations) > 0 {
		return locations
	}

	return nil
}

// Generates random mine locations, excluding 3x3 square around provided coordinates.
func (g *Game) randomMineLocations(aroundX, aroundY int) []int {
	locations := make([]int, 0, g.minesToPlant)
	for _, i := range rand.Perm(len(g.cells)) {
		if len(locations) == g.minesToPlant {
			break
		}

		x := i % g.width
		y := i / g.width
		if x < aroundX-1 || x > aroundX+1 || y < aroundY-1 || y > aroundY+1 {
			locations = append(locations, i)
		}
	}

	return locations
}

// Plants mines into specified locations.
func (g *Game) plantMines(mineLocations []int) {
	// Set up mines and make sure we don't count them twice
	for _, i := range mineLocations {
		if !g.cells[i].isMine {
			g.cells[i].isMine = true
			g.minesLeft++
		}
	}

	// Pre-calculate all adjacent numbers
	empty := 0
	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			cell := g.Cell(x, y)
			cell.adjacentMines = 0
			for _, offset := range adjacentOffsets {
				if g.Cell(x+offset.x, y+offset.y).isMine {
					cell.adjacentMines++
				}
			}
			if cell.adjacentMines == 0 {
				empty++
			}
		}
	}

	// Calculate heart spawn threshold
	g.heartSpawnThreshold = empty / (g.heartsLeft + 1)
	if g.heartSpawnThreshold == 0 {
		g.heartSpawnThreshold = 1
	}
}
