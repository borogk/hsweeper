package ui

import (
	"minesweeper/game"

	"github.com/gdamore/tcell/v2"
)

type TitleMenuView struct {
	ui *Ui
}

var logo = [][]rune{
	[]rune("██   ██  ░░░░░░  ░░      ░░  ░░░░░░  ░░░░░░  ░░░░░░  ░░░░░░  ░░░░  "),
	[]rune("██   ██  ░░      ░░      ░░  ░░      ░░      ░░  ░░  ░░      ░░  ░░"),
	[]rune("███████  ░░░░░░  ░░  ░░  ░░  ░░░░    ░░░░    ░░░░░░  ░░░░    ░░  ░░"),
	[]rune("██   ██      ░░  ░░  ░░  ░░  ░░      ░░      ░░      ░░      ░░░░  "),
	[]rune("██   ██  ░░░░░░    ░░  ░░    ░░░░░░  ░░░░░░  ░░      ░░░░░░  ░░  ░░"),
}

func newDifficultySelectView(ui *Ui) *TitleMenuView {
	return &TitleMenuView{ui: ui}
}

func (v *TitleMenuView) ContentSize() (width, height int) {
	return len(logo[0]), len(logo) + 10
}

func (v *TitleMenuView) Draw(screen tcell.Screen) {
	screenWidth, screenHeight := screen.Size()
	contentWidth, contentHeight := v.ContentSize()
	palette := defaultPalette

	logoX := (screenWidth - contentWidth) / 2
	logoY := (screenHeight - contentHeight) / 2
	for y, logoLine := range logo {
		for x, c := range logoLine {
			logoStyle := palette.Blank
			if c == '█' {
				logoStyle = palette.Logo
			} else if c == '░' {
				logoStyle = palette.LogoSecondary
			}

			screen.Put(logoX+x, logoY+y, " ", logoStyle)
		}
	}

	hint := "(select the option)"
	hintX := (screenWidth - len(hint)) / 2
	hintY := logoY + len(logo) + 2
	screen.PutStrStyled(hintX, hintY, hint, palette.HintText)

	screen.PutStrStyled(hintX+1, hintY+2, "[1]  H-Expert", palette.ExpertGameText)
	screen.PutStrStyled(hintX+1, hintY+3, "[2]  H-Big", palette.BigGameText)

	screen.PutStrStyled(hintX+1, hintY+5, "[3]  Classic Easy", palette.ClassicGameText)
	screen.PutStrStyled(hintX+1, hintY+6, "[4]  Classic Medium", palette.ClassicGameText)
	screen.PutStrStyled(hintX+1, hintY+7, "[5]  Classic Expert", palette.ClassicGameText)
}

func (v *TitleMenuView) OnInput(_ tcell.Key, rune rune) {
	switch rune {
	case '1':
		v.ui.pushView(newGameView(v.ui, newExpertGame()))
	case '2':
		v.ui.pushView(newGameView(v.ui, v.newBigGame()))
	case '3':
		v.ui.pushView(newGameView(v.ui, newClassicGame(9, 9, 10)))
	case '4':
		v.ui.pushView(newGameView(v.ui, newClassicGame(16, 16, 40)))
	case '5':
		v.ui.pushView(newGameView(v.ui, newClassicGame(30, 16, 99)))
	}
}

func newExpertGame() *game.Game {
	return game.NewGame(30, 16, 99, 1, 1)
}

func (v *TitleMenuView) newBigGame() *game.Game {
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
	hearts := cells/480 - cells/2400
	extraLives := cells / 2400
	return game.NewGame(gameWidth, gameHeight, mines, hearts, 1+extraLives)
}

func newClassicGame(width, height, mines int) *game.Game {
	return game.NewGame(width, height, mines, 0, 1)
}
