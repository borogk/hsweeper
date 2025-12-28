package game

// Cell represents a single game cell
type Cell struct {
	isMine        bool
	isHeart       bool
	isRevealed    bool
	isFlagged     bool
	isQuestioned  bool
	adjacentMines int
}

// IsMine indicates if the cell has a mine planted
func (c *Cell) IsMine() bool {
	return c.isMine
}

// IsHeart indicates if the cell has a heart pickup
func (c *Cell) IsHeart() bool {
	return c.isHeart
}

// IsRevealed indicates if the cell has been revealed
func (c *Cell) IsRevealed() bool {
	return c.isRevealed
}

// IsFlagged indicates if the cell was marked with a flag
func (c *Cell) IsFlagged() bool {
	return c.isFlagged
}

// IsQuestioned indicates if the cell was marked with a question
func (c *Cell) IsQuestioned() bool {
	return c.isQuestioned
}

// AdjacentMines returns precalculated amount of adjacent mines
func (c *Cell) AdjacentMines() int {
	return c.adjacentMines
}
