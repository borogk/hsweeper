package ui

import (
	"hsweeper/game"

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

	optionsX := logoX + 4
	optionsY := logoY + len(logo) + 2

	screen.PutStrStyled(optionsX, optionsY, "[1] H-Expert", palette.ExpertGameText)
	screen.PutStrStyled(optionsX, optionsY+1, "[2] H-Big", palette.BigGameText)
	screen.PutStrStyled(optionsX, optionsY+2, "[3] Classic Easy", palette.ClassicGameText)
	screen.PutStrStyled(optionsX, optionsY+3, "[4] Classic Medium", palette.ClassicGameText)
	screen.PutStrStyled(optionsX, optionsY+4, "[5] Classic Expert", palette.ClassicGameText)

	hint := "(press SPACE to start, 1-5 to select an option, ESC to quit)"
	hintX := logoX + 2
	hintY := optionsY + 6
	screen.PutStrStyled(hintX, hintY, hint, palette.HintText)
}

func (v *TitleMenuView) OnInput(key tcell.Key, rune rune) {
	switch key {
	case tcell.KeyEscape:
		v.ui.popView()
	default:
		switch rune {
		case '1', ' ':
			v.startGame(newExpertGame())
		case '2':
			v.startGame(v.newBigGame())
		case '3':
			v.startGame(newClassicGame(9, 9, 10))
		case '4':
			v.startGame(newClassicGame(16, 16, 40))
		case '5':
			v.startGame(newClassicGame(30, 16, 99))
		}
	}
}

func (v *TitleMenuView) startGame(gameFactory func() *game.Game) {
	v.ui.pushView(newGameView(v.ui, gameFactory))
}

func newExpertGame() func() *game.Game {
	return func() *game.Game {
		return game.NewGame(30, 16, 99, 1, 1)
	}
}

func (v *TitleMenuView) newBigGame() func() *game.Game {
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

func newClassicGame(width, height, mines int) func() *game.Game {
	return func() *game.Game {
		return game.NewGame(width, height, mines, 0, 1)
	}
}
