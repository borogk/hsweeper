package game

type Point struct {
	x int
	y int
}

// Helper offsets to quickly determine adjacent cells.
var adjacentOffsets = []Point{
	{x: -1, y: -1},
	{x: 0, y: -1},
	{x: 1, y: -1},
	{x: 1, y: 0},
	{x: 1, y: 1},
	{x: 0, y: 1},
	{x: -1, y: 1},
	{x: -1, y: 0},
}
