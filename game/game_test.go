package game

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	t.Run("creates ready and empty game", func(t *testing.T) {
		g := NewGame(9, 8, 7, 6, 5)

		assertEquals(t, g.status, StatusReady)
		assertEquals(t, g.width, 9)
		assertEquals(t, g.height, 8)
		assertEquals(t, g.minesToPlant, 7)
		assertEquals(t, g.heartsToPlant, 6)
		assertEquals(t, g.heartsLeft, 6)
		assertEquals(t, g.livesLeft, 5)
		assertEquals(t, g.unrevealedCounter, 72)
		assertEquals(t, g.minesLeft, 0)
		assertEquals(t, g.flaggedCounter, 0)
		assertEquals(t, g.heartSpawnCounter, 0)
		assertEquals(t, g.heartSpawnThreshold, 0)

		assertEquals(t, len(g.cells), 72)
		for _, c := range g.cells {
			assertEquals(t, c.isMine, false)
			assertEquals(t, c.isHeart, false)
			assertEquals(t, c.isRevealed, false)
			assertEquals(t, c.isFlagged, false)
			assertEquals(t, c.isQuestioned, false)
			assertEquals(t, c.adjacentMines, 0)
		}
	})

	t.Run("adjusts to minimum width of 1", func(t *testing.T) {
		g := NewGame(-1, 8, 7, 6, 5)
		assertEquals(t, g.width, 1)
	})

	t.Run("adjusts to minimum height of 1", func(t *testing.T) {
		g := NewGame(9, -1, 7, 6, 5)
		assertEquals(t, g.height, 1)
	})

	t.Run("adjusts to minimum mines of 0", func(t *testing.T) {
		g := NewGame(9, 8, -1, 6, 5)
		assertEquals(t, g.minesToPlant, 0)
	})

	t.Run("adjusts to minimum hearts of 0", func(t *testing.T) {
		g := NewGame(9, 8, 7, -1, 5)
		assertEquals(t, g.heartsToPlant, 0)
		assertEquals(t, g.heartsLeft, 0)
	})

	t.Run("adjusts to minimum lives of 0 and sets status to lost", func(t *testing.T) {
		g := NewGame(9, 8, 7, 6, -1)
		assertEquals(t, g.livesLeft, 0)
		assertEquals(t, g.status, StatusLost)
	})
}

func TestGame_Status(t *testing.T) {
	g := &Game{status: StatusStarted}
	assertEquals(t, g.Status(), StatusStarted)
}

func TestGame_Width(t *testing.T) {
	g := &Game{width: 100}
	assertEquals(t, g.Width(), 100)
}

func TestGame_Height(t *testing.T) {
	g := &Game{height: 200}
	assertEquals(t, g.Height(), 200)
}

func TestGame_IsOutOfBounds(t *testing.T) {
	g := NewGame(2, 3, 0, 0, 1)

	t.Run("returns true for all inner cells", func(t *testing.T) {
		assertEquals(t, g.IsOutOfBounds(0, 0), false)
		assertEquals(t, g.IsOutOfBounds(1, 0), false)
		assertEquals(t, g.IsOutOfBounds(0, 1), false)
		assertEquals(t, g.IsOutOfBounds(1, 1), false)
		assertEquals(t, g.IsOutOfBounds(0, 2), false)
		assertEquals(t, g.IsOutOfBounds(1, 2), false)
	})

	t.Run("returns true for outer cells", func(t *testing.T) {
		assertEquals(t, g.IsOutOfBounds(-1, 0), true)
		assertEquals(t, g.IsOutOfBounds(2, 0), true)
		assertEquals(t, g.IsOutOfBounds(0, -1), true)
		assertEquals(t, g.IsOutOfBounds(0, 3), true)
		assertEquals(t, g.IsOutOfBounds(-1, -1), true)
		assertEquals(t, g.IsOutOfBounds(3, 3), true)
	})
}

func TestGame_Cell(t *testing.T) {
	g := NewGame(2, 3, 0, 0, 1)

	t.Run("returns inner cells", func(t *testing.T) {
		assertSame(t, g.Cell(0, 0), &g.cells[0])
		assertSame(t, g.Cell(1, 0), &g.cells[1])
		assertSame(t, g.Cell(0, 1), &g.cells[2])
		assertSame(t, g.Cell(1, 1), &g.cells[3])
		assertSame(t, g.Cell(0, 2), &g.cells[4])
		assertSame(t, g.Cell(1, 2), &g.cells[5])
	})

	t.Run("returns fake outer cell", func(t *testing.T) {
		c := g.Cell(-1, -1)

		assertEquals(t, c.isRevealed, true)
		assertNotSame(t, c, &g.cells[0])
		assertNotSame(t, c, &g.cells[1])
		assertNotSame(t, c, &g.cells[2])
		assertNotSame(t, c, &g.cells[3])
		assertNotSame(t, c, &g.cells[4])
		assertNotSame(t, c, &g.cells[5])
	})
}

func TestGame_IsFinished(t *testing.T) {
	assertEquals(t, (&Game{status: StatusReady}).IsFinished(), false)
	assertEquals(t, (&Game{status: StatusStarted}).IsFinished(), false)
	assertEquals(t, (&Game{status: StatusLost}).IsFinished(), true)
	assertEquals(t, (&Game{status: StatusWon}).IsFinished(), true)
}

func TestGame_MinesRemaining(t *testing.T) {
	baseSnapshot := Snapshot{
		Width:        3,
		Height:       3,
		MinesToPlant: 2,
		LivesLeft:    1,
	}

	t.Run("returns 0 before game starts", func(t *testing.T) {
		snapshot := baseSnapshot
		snapshot.Status = StatusReady
		g := RestoreGameFromSnapshot(snapshot)

		assertEquals(t, g.MinesRemaining(), 0)
	})

	t.Run("returns all planted mines after game starts", func(t *testing.T) {
		snapshot := baseSnapshot
		snapshot.Status = StatusStarted
		snapshot.MineLocations = locationsFromBitmap(
			"---",
			"---",
			"xxx",
		)
		g := RestoreGameFromSnapshot(snapshot)

		assertEquals(t, g.MinesRemaining(), 3)
	})

	t.Run("returns less mines after adding flags", func(t *testing.T) {
		snapshot := baseSnapshot
		snapshot.Status = StatusStarted
		snapshot.MineLocations = locationsFromBitmap(
			"---",
			"---",
			"xxx",
		)
		g := RestoreGameFromSnapshot(snapshot)
		g.ToggleFlag(1, 2)
		g.ToggleFlag(2, 2)

		assertEquals(t, g.MinesRemaining(), 1)
	})

}

func TestGame_LivesRemaining(t *testing.T) {
	g := &Game{livesLeft: 3}
	assertEquals(t, g.LivesRemaining(), 3)
}

func TestGame_ToggleFlag(t *testing.T) {
	toggle := func(g *Game, x, y int) { g.ToggleFlag(x, y) }
	toggleOther := func(g *Game, x, y int) { g.ToggleQuestion(x, y) }
	checkMark := isCellFlagged
	checkOtherMark := isCellQuestioned

	testToggleMark(t, toggle, toggleOther, checkMark, checkOtherMark)
}

func TestGame_ToggleQuestion(t *testing.T) {
	toggle := func(g *Game, x, y int) { g.ToggleQuestion(x, y) }
	toggleOther := func(g *Game, x, y int) { g.ToggleFlag(x, y) }
	checkMark := isCellQuestioned
	checkOtherMark := isCellFlagged

	testToggleMark(t, toggle, toggleOther, checkMark, checkOtherMark)
}

func testToggleMark(
	t *testing.T,
	toggle func(g *Game, x, y int),
	toggleOther func(g *Game, x, y int),
	checkMark func(c *Cell) bool,
	checkOtherMark func(c *Cell) bool,
) {
	snapshot := Snapshot{
		Status:    StatusStarted,
		Width:     3,
		Height:    3,
		LivesLeft: 1,
	}

	t.Run("adds and removes mark", func(t *testing.T) {
		g := RestoreGameFromSnapshot(snapshot)

		beforeToggle := g.toBitmap(checkMark)
		toggle(g, 1, 2)
		afterToggledFirst := g.toBitmap(checkMark)
		toggle(g, 2, 2)
		afterToggledSecond := g.toBitmap(checkMark)
		toggle(g, 1, 2)
		afterToggledFirstTwice := g.toBitmap(checkMark)
		toggle(g, 2, 2)
		afterToggledSecondTwice := g.toBitmap(checkMark)

		assertBitmapEquals(t, beforeToggle,
			"---",
			"---",
			"---",
		)
		assertBitmapEquals(t, afterToggledFirst,
			"---",
			"---",
			"-x-",
		)
		assertBitmapEquals(t, afterToggledSecond,
			"---",
			"---",
			"-xx",
		)
		assertBitmapEquals(t, afterToggledFirstTwice,
			"---",
			"---",
			"--x",
		)
		assertBitmapEquals(t, afterToggledSecondTwice,
			"---",
			"---",
			"---",
		)
	})

	t.Run("replaces another mark", func(t *testing.T) {
		g := RestoreGameFromSnapshot(snapshot)
		toggleOther(g, 1, 1)

		toggle(g, 1, 1)

		cell := g.Cell(1, 1)
		assertEquals(t, checkMark(cell), true)
		assertEquals(t, checkOtherMark(cell), false)
	})
}

// TODO: test other functions
