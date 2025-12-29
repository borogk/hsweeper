package game

type CellPredicate func(*Cell) bool

func isCellMine(c *Cell) bool {
	return c.isMine
}

func isCellHeart(c *Cell) bool {
	return c.isHeart
}

func isCellRevealed(c *Cell) bool {
	return c.isRevealed
}

func isCellFlagged(c *Cell) bool {
	return c.isFlagged
}

func isCellQuestioned(c *Cell) bool {
	return c.isQuestioned
}
