package ui

import (
	"github.com/borogk/hsweeper/game"
)

type GameFactory func() *game.Game

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

func newExpertGameFactory() GameFactory {
	return func() *game.Game {
		return game.NewGame(30, 16, 99, 1, 1)
	}
}

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

func newClassicGameFactory(width, height, mines int) GameFactory {
	return func() *game.Game {
		return game.NewGame(width, height, mines, 0, 1)
	}
}
