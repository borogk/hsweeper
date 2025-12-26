package game

type Cell struct {
	isMine        bool
	isHeart       bool
	isRevealed    bool
	isFlagged     bool
	isQuestioned  bool
	adjacentMines int
}

func (c *Cell) IsMine() bool {
	return c.isMine
}

func (c *Cell) IsHeart() bool {
	return c.isHeart
}

func (c *Cell) IsRevealed() bool {
	return c.isRevealed
}

func (c *Cell) IsFlagged() bool {
	return c.isFlagged
}

func (c *Cell) IsQuestioned() bool {
	return c.isQuestioned
}

func (c *Cell) AdjacentMines() int {
	return c.adjacentMines
}
