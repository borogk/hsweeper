package game

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func assertEquals(t *testing.T, actual, expected any) {
	t.Helper()
	diff := cmp.Diff(actual, expected)
	if diff != "" {
		t.Errorf("assertEquals fails\n%s", diff)
	}
}

func assertSame(t *testing.T, actual, expected any) {
	t.Helper()
	if actual != expected {
		t.Errorf("assertSame fails")
	}
}

func assertNotSame(t *testing.T, actual, expected any) {
	t.Helper()
	if actual == expected {
		t.Errorf("assertNotSame fails")
	}
}

func assertBitmapEquals(t *testing.T, actual []string, expected ...string) {
	t.Helper()
	assertEquals(t, actual, expected)
}

func (g *Game) toBitmap(bitmapFunc func(c *Cell) bool) []string {
	result := make([]string, g.height)
	for y := 0; y < g.height; y++ {
		line := make([]rune, g.width)
		for x := 0; x < g.width; x++ {
			if bitmapFunc(g.Cell(x, y)) {
				line[x] = 'x'
			} else {
				line[x] = '-'
			}
		}
		result[y] = string(line)
	}
	return result
}

func (g *Game) toNumbersMap() []string {
	result := make([]string, g.height)
	for y := 0; y < g.height; y++ {
		line := make([]rune, g.width)
		for x := 0; x < g.width; x++ {
			adjacentMines := g.Cell(x, y).adjacentMines
			if adjacentMines > 0 {
				line[x] = '0' + rune(adjacentMines)
			} else {
				line[x] = '-'
			}
		}
		result[y] = string(line)
	}
	return result
}

func locationsFromBitmap(bitmap ...string) []int {
	locations := make([]int, 0)
	for y, line := range bitmap {
		for x, r := range line {
			if r == 'x' {
				locations = append(locations, x+y*len(line))
			}
		}
	}
	return locations
}

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
