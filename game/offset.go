package game

type Offset struct {
	x int
	y int
}

// Helper offsets to quickly determine adjacent cells.
var adjacentOffsets = []Offset{
	{x: -1, y: -1},
	{x: 0, y: -1},
	{x: 1, y: -1},
	{x: 1, y: 0},
	{x: 1, y: 1},
	{x: 0, y: 1},
	{x: -1, y: 1},
	{x: -1, y: 0},
}
