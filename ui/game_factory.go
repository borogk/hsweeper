package ui

import (
	"github.com/borogk/hsweeper/game"
)

// GameFactory repeatedly creates new game instances to facilitate restarts.
// May return nil, which means restarting the game is impossible.
type GameFactory func() *game.Game

// Special game factory, that returns an existing game only once.
func newExistingGameFactory(g *game.Game) GameFactory {
	once := true
	return func() *game.Game {
		if once {
			once = false
			return g
		}

		return nil
	}
}

// H-Expert game factory.
func newExpertGameFactory() GameFactory {
	return func() *game.Game {
		return game.NewGame(30, 16, 99, 1, 1)
	}
}

// H-Big game factory.
func (v *TitleMenuView) newBigGameFactory() GameFactory {
	return func() *game.Game {
		width, height := v.ui.screen.Size()
		gameWidth := (width - 2) / 3
		gameHeight := height - 5
		if gameWidth < 30 {
			gameWidth = 30
		}
		if gameHeight < 16 {
			gameHeight = 16
		}

		cells := gameWidth * gameHeight
		mines := cells / 5
		if mines < 99 {
			mines = 99
		}
		hearts := cells/480 - cells/2400
		extraLives := cells / 2400
		return game.NewGame(gameWidth, gameHeight, mines, hearts, 1+extraLives)
	}
}

// Classic game factory.
func newClassicGameFactory(width, height, mines int) GameFactory {
	return func() *game.Game {
		return game.NewGame(width, height, mines, 0, 1)
	}
}
