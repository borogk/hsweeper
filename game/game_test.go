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

func TestGame_RestoreSave(t *testing.T) {
	snapshot := &Snapshot{
		Status:        StatusStarted,
		Width:         6,
		Height:        5,
		MinesToPlant:  8,
		HeartsToPlant: 1,
		LivesLeft:     4,
		HeartsLeft:    1,
		MineLocations: locationsFromBitmap(
			"------",
			"--xx--",
			"-xxxx-",
			"--xx--",
			"------",
		),
		RevealedLocations: locationsFromBitmap(
			"xxxxxx",
			"x----x",
			"x----x",
			"x----x",
			"xxxxxx",
		),
		UncollectedHeartLocations: locationsFromBitmap(
			"x----x",
			"------",
			"------",
			"------",
			"------",
		),
		FlaggedLocations: locationsFromBitmap(
			"------",
			"--xx--",
			"------",
			"------",
			"------",
		),
		QuestionedLocations: locationsFromBitmap(
			"------",
			"------",
			"------",
			"--xx--",
			"------",
		),
	}

	t.Run("RestoreGame correctly restores all game aspects", func(t *testing.T) {
		g := RestoreGame(snapshot)

		assertEquals(t, g.status, StatusStarted)
		assertEquals(t, g.width, 6)
		assertEquals(t, g.height, 5)
		assertEquals(t, g.minesToPlant, 8)
		assertEquals(t, g.heartsToPlant, 1)
		assertEquals(t, g.livesLeft, 4)
		assertEquals(t, g.heartsLeft, 1)
		assertEquals(t, g.unrevealedCounter, 12)
		assertEquals(t, g.flaggedCounter, 2)
		assertEquals(t, g.heartSpawnCounter, 4)
		assertEquals(t, g.heartSpawnThreshold, 2)

		assertBitmapEquals(t, g.toNumbersMap(),
			"012210",
			"134431",
			"136631",
			"134431",
			"012210",
		)
		assertBitmapEquals(t, g.toBitmap(isCellMine),
			"------",
			"--xx--",
			"-xxxx-",
			"--xx--",
			"------",
		)
		assertBitmapEquals(t, g.toBitmap(isCellRevealed),
			"xxxxxx",
			"x----x",
			"x----x",
			"x----x",
			"xxxxxx",
		)
		assertBitmapEquals(t, g.toBitmap(isCellHeart),
			"x----x",
			"------",
			"------",
			"------",
			"------",
		)
		assertBitmapEquals(t, g.toBitmap(isCellFlagged),
			"------",
			"--xx--",
			"------",
			"------",
			"------",
		)
		assertBitmapEquals(t, g.toBitmap(isCellQuestioned),
			"------",
			"------",
			"------",
			"--xx--",
			"------",
		)
	})

	t.Run("Save creates correct snapshot", func(t *testing.T) {
		save := RestoreGame(snapshot).Save()

		assertEquals(t, save.Status, StatusStarted)
		assertEquals(t, save.Width, 6)
		assertEquals(t, save.Height, 5)
		assertEquals(t, save.MinesToPlant, 8)
		assertEquals(t, save.HeartsToPlant, 1)
		assertEquals(t, save.LivesLeft, 4)
		assertEquals(t, save.HeartsLeft, 1)
		assertEquals(t, save.MineLocations, []int{8, 9, 13, 14, 15, 16, 20, 21})
		assertEquals(t, save.RevealedLocations, []int{0, 1, 2, 3, 4, 5, 6, 11, 12, 17, 18, 23, 24, 25, 26, 27, 28, 29})
		assertEquals(t, save.UncollectedHeartLocations, []int{0, 5})
		assertEquals(t, save.FlaggedLocations, []int{8, 9})
		assertEquals(t, save.QuestionedLocations, []int{20, 21})
	})

	t.Run("save-restore-save again produce two identical snapshots", func(t *testing.T) {
		firstSave := RestoreGame(snapshot).Save()
		secondSave := RestoreGame(firstSave).Save()
		assertEquals(t, firstSave, secondSave)
	})

	t.Run("RestoreGame doesn't reveal twice", func(t *testing.T) {
		g := RestoreGame(&Snapshot{
			Status:            StatusStarted,
			Width:             3,
			Height:            3,
			LivesLeft:         1,
			RevealedLocations: []int{0, 1, 0, 1},
		})

		assertEquals(t, g.unrevealedCounter, 7)
		assertBitmapEquals(t, g.toBitmap(isCellRevealed),
			"xx-",
			"---",
			"---",
		)
	})

	t.Run("RestoreGame doesn't reveal mined cells", func(t *testing.T) {
		g := RestoreGame(&Snapshot{
			Status:            StatusStarted,
			Width:             3,
			Height:            3,
			LivesLeft:         1,
			MineLocations:     []int{0, 1},
			RevealedLocations: []int{0, 1, 2},
		})

		assertEquals(t, g.unrevealedCounter, 8)
		assertBitmapEquals(t, g.toBitmap(isCellRevealed),
			"--x",
			"---",
			"---",
		)
	})

	t.Run("RestoreGame doesn't put hearts twice", func(t *testing.T) {
		g := RestoreGame(&Snapshot{
			Status:                    StatusStarted,
			Width:                     3,
			Height:                    3,
			LivesLeft:                 1,
			RevealedLocations:         []int{0, 1},
			UncollectedHeartLocations: []int{0, 1, 0, 1},
		})

		assertBitmapEquals(t, g.toBitmap(isCellHeart),
			"xx-",
			"---",
			"---",
		)
	})

	t.Run("RestoreGame doesn't put hearts on unrevealed cells", func(t *testing.T) {
		g := RestoreGame(&Snapshot{
			Status:                    StatusStarted,
			Width:                     3,
			Height:                    3,
			LivesLeft:                 1,
			RevealedLocations:         []int{0, 1},
			UncollectedHeartLocations: []int{0, 1, 2},
		})

		assertBitmapEquals(t, g.toBitmap(isCellHeart),
			"xx-",
			"---",
			"---",
		)
	})

	t.Run("RestoreGame doesn't put hearts on bomb-adjacent cells", func(t *testing.T) {
		g := RestoreGame(&Snapshot{
			Status:    StatusStarted,
			Width:     3,
			Height:    3,
			LivesLeft: 1,
			MineLocations: locationsFromBitmap(
				"---",
				"---",
				"--x",
			),
			RevealedLocations: locationsFromBitmap(
				"xxx",
				"xxx",
				"---",
			),
			UncollectedHeartLocations: locationsFromBitmap(
				"xxx",
				"xxx",
				"---",
			),
		})

		assertBitmapEquals(t, g.toBitmap(isCellHeart),
			"xxx",
			"x--",
			"---",
		)
	})

	t.Run("RestoreGame doesn't put flags twice", func(t *testing.T) {
		g := RestoreGame(&Snapshot{
			Status:           StatusStarted,
			Width:            3,
			Height:           3,
			LivesLeft:        1,
			FlaggedLocations: []int{0, 1, 0, 1},
		})

		assertEquals(t, g.flaggedCounter, 2)
		assertBitmapEquals(t, g.toBitmap(isCellFlagged),
			"xx-",
			"---",
			"---",
		)
	})

	t.Run("RestoreGame doesn't put flags on revealed cells", func(t *testing.T) {
		g := RestoreGame(&Snapshot{
			Status:            StatusStarted,
			Width:             3,
			Height:            3,
			LivesLeft:         1,
			RevealedLocations: []int{0, 1},
			FlaggedLocations:  []int{0, 1, 2},
		})

		assertEquals(t, g.flaggedCounter, 1)
		assertBitmapEquals(t, g.toBitmap(isCellFlagged),
			"--x",
			"---",
			"---",
		)
	})

	t.Run("RestoreGame doesn't put questions on revealed cells", func(t *testing.T) {
		g := RestoreGame(&Snapshot{
			Status:              StatusStarted,
			Width:               3,
			Height:              3,
			LivesLeft:           1,
			RevealedLocations:   []int{0, 1},
			QuestionedLocations: []int{0, 1, 2},
		})

		assertBitmapEquals(t, g.toBitmap(isCellQuestioned),
			"--x",
			"---",
			"---",
		)
	})

	t.Run("RestoreGame doesn't put questions on flags", func(t *testing.T) {
		g := RestoreGame(&Snapshot{
			Status:              StatusStarted,
			Width:               3,
			Height:              3,
			LivesLeft:           1,
			FlaggedLocations:    []int{0, 1},
			QuestionedLocations: []int{0, 1, 2},
		})

		assertEquals(t, g.flaggedCounter, 2)
		assertBitmapEquals(t, g.toBitmap(isCellFlagged),
			"xx-",
			"---",
			"---",
		)
		assertBitmapEquals(t, g.toBitmap(isCellQuestioned),
			"--x",
			"---",
			"---",
		)
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
	baseSnapshot := &Snapshot{
		Width:        3,
		Height:       3,
		MinesToPlant: 2,
		LivesLeft:    1,
	}

	t.Run("returns 0 before game starts", func(t *testing.T) {
		snapshot := baseSnapshot
		snapshot.Status = StatusReady
		g := RestoreGame(snapshot)

		assertEquals(t, g.MinesRemaining(), 2)
	})

	t.Run("returns all planted mines after game starts", func(t *testing.T) {
		snapshot := baseSnapshot
		snapshot.Status = StatusStarted
		snapshot.MineLocations = locationsFromBitmap(
			"---",
			"---",
			"xxx",
		)
		g := RestoreGame(snapshot)

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
		g := RestoreGame(snapshot)
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
	checkMark CellPredicate,
	checkOtherMark CellPredicate,
) {
	snapshot := &Snapshot{
		Status:    StatusStarted,
		Width:     3,
		Height:    3,
		LivesLeft: 1,
	}

	t.Run("adds and removes mark", func(t *testing.T) {
		g := RestoreGame(snapshot)

		toggle(g, 1, 2)
		afterToggledFirst := g.toBitmap(checkMark)
		toggle(g, 2, 2)
		afterToggledSecond := g.toBitmap(checkMark)
		toggle(g, 1, 2)
		afterToggledFirstTwice := g.toBitmap(checkMark)
		toggle(g, 2, 2)
		afterToggledSecondTwice := g.toBitmap(checkMark)

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
		g := RestoreGame(snapshot)
		toggleOther(g, 1, 1)

		toggle(g, 1, 1)

		cell := g.Cell(1, 1)
		assertEquals(t, checkMark(cell), true)
		assertEquals(t, checkOtherMark(cell), false)
	})

	for _, status := range []Status{StatusLost, StatusWon} {
		t.Run("does nothing of finished game", func(t *testing.T) {
			g := RestoreGame(snapshot)
			g.status = status

			for x := 0; x < 3; x++ {
				for y := 0; y < 3; y++ {
					toggle(g, x, y)
				}
			}
			afterToggledAll := g.toBitmap(checkMark)

			assertBitmapEquals(t, afterToggledAll,
				"---",
				"---",
				"---",
			)
		})
	}
}

func TestGame_ClearFlagAndQuestion(t *testing.T) {
	snapshot := &Snapshot{
		Status:    StatusStarted,
		Width:     3,
		Height:    3,
		LivesLeft: 1,
		FlaggedLocations: locationsFromBitmap(
			"xx-",
			"---",
			"---",
		),
		QuestionedLocations: locationsFromBitmap(
			"---",
			"---",
			"-xx",
		),
	}

	t.Run("clears flags", func(t *testing.T) {
		g := RestoreGame(snapshot)

		g.ClearFlagAndQuestion(0, 0)
		afterClearedOne := g.toBitmap(isCellFlagged)
		g.ClearFlagAndQuestion(1, 0)
		afterClearedBoth := g.toBitmap(isCellFlagged)

		assertBitmapEquals(t, afterClearedOne,
			"-x-",
			"---",
			"---",
		)
		assertBitmapEquals(t, afterClearedBoth,
			"---",
			"---",
			"---",
		)
	})

	t.Run("clears questions", func(t *testing.T) {
		g := RestoreGame(snapshot)

		g.ClearFlagAndQuestion(1, 2)
		afterClearedOne := g.toBitmap(isCellQuestioned)
		g.ClearFlagAndQuestion(2, 2)
		afterClearedBoth := g.toBitmap(isCellQuestioned)

		assertBitmapEquals(t, afterClearedOne,
			"---",
			"---",
			"--x",
		)
		assertBitmapEquals(t, afterClearedBoth,
			"---",
			"---",
			"---",
		)
	})

	t.Run("does nothing on unmarked cells", func(t *testing.T) {
		g := RestoreGame(snapshot)

		g.ClearFlagAndQuestion(1, 1)
		flagsAfterClear := g.toBitmap(isCellFlagged)
		questionsAfterClear := g.toBitmap(isCellQuestioned)

		assertBitmapEquals(t, flagsAfterClear,
			"xx-",
			"---",
			"---",
		)
		assertBitmapEquals(t, questionsAfterClear,
			"---",
			"---",
			"-xx",
		)
	})

	for _, status := range []Status{StatusLost, StatusWon} {
		t.Run("does nothing on finished game", func(t *testing.T) {
			g := RestoreGame(snapshot)
			g.status = status

			for x := 0; x < 3; x++ {
				for y := 0; y < 3; y++ {
					g.ClearFlagAndQuestion(x, y)
				}
			}
			flagsAfterClearedAll := g.toBitmap(isCellFlagged)
			questionsAfterClearedAll := g.toBitmap(isCellQuestioned)

			assertBitmapEquals(t, flagsAfterClearedAll,
				"xx-",
				"---",
				"---",
			)
			assertBitmapEquals(t, questionsAfterClearedAll,
				"---",
				"---",
				"-xx",
			)
		})
	}
}

func TestGame_Pickup(t *testing.T) {
	snapshot := &Snapshot{
		Status:    StatusStarted,
		Width:     3,
		Height:    3,
		LivesLeft: 1,
		RevealedLocations: locationsFromBitmap(
			"xx-",
			"---",
			"---",
		),
		UncollectedHeartLocations: locationsFromBitmap(
			"xx-",
			"---",
			"---",
		),
	}

	t.Run("picks up extra lives", func(t *testing.T) {
		g := RestoreGame(snapshot)

		g.Pickup(0, 0)
		heartsAfterPickedOne := g.toBitmap(isCellHeart)
		livesAfterPickedOne := g.livesLeft
		g.Pickup(1, 0)
		heartsAfterPickedBoth := g.toBitmap(isCellHeart)
		livesAfterPickedBoth := g.livesLeft

		assertBitmapEquals(t, heartsAfterPickedOne,
			"-x-",
			"---",
			"---",
		)
		assertBitmapEquals(t, heartsAfterPickedBoth,
			"---",
			"---",
			"---",
		)
		assertEquals(t, livesAfterPickedOne, 2)
		assertEquals(t, livesAfterPickedBoth, 3)
	})

	t.Run("does nothing on cells without hearts", func(t *testing.T) {
		g := RestoreGame(snapshot)

		g.Pickup(1, 1)
		heartsAfterPicked := g.toBitmap(isCellHeart)
		livesAfterPicked := g.livesLeft

		assertBitmapEquals(t, heartsAfterPicked,
			"xx-",
			"---",
			"---",
		)
		assertEquals(t, livesAfterPicked, 1)
	})

	for _, status := range []Status{StatusLost, StatusWon} {
		t.Run("does nothing on finished game", func(t *testing.T) {
			g := RestoreGame(snapshot)
			g.status = status

			for x := 0; x < 3; x++ {
				for y := 0; y < 3; y++ {
					g.Pickup(x, y)
				}
			}
			heartsAfterPickedAll := g.toBitmap(isCellHeart)
			livesAfterPickedAll := g.livesLeft

			assertBitmapEquals(t, heartsAfterPickedAll,
				"xx-",
				"---",
				"---",
			)
			assertEquals(t, livesAfterPickedAll, 1)
		})
	}
}

// TODO: test other functions
