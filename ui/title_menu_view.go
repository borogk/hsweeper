package ui

import (
	"fmt"

	"github.com/borogk/hsweeper/game"
	"github.com/gdamore/tcell/v2"
)

type TitleMenuView struct {
	ui        *Ui
	savedGame *game.Game
}

var logo = [][]rune{
	[]rune("██   ██  ░░░░░░  ░░      ░░  ░░░░░░  ░░░░░░  ░░░░░░  ░░░░░░  ░░░░  "),
	[]rune("██   ██  ░░      ░░      ░░  ░░      ░░      ░░  ░░  ░░      ░░  ░░"),
	[]rune("███████  ░░░░░░  ░░  ░░  ░░  ░░░░    ░░░░    ░░░░░░  ░░░░    ░░  ░░"),
	[]rune("██   ██      ░░  ░░  ░░  ░░  ░░      ░░      ░░      ░░      ░░░░  "),
	[]rune("██   ██  ░░░░░░    ░░  ░░    ░░░░░░  ░░░░░░  ░░      ░░░░░░  ░░  ░░"),
}

func newTitleMenuView(ui *Ui) *TitleMenuView {
	return &TitleMenuView{ui: ui}
}

func (v *TitleMenuView) OnActivate() {
	v.savedGame = game.LoadGame(game.DefaultSavePath())
}

func (v *TitleMenuView) OnDeactivate() {

}

func (v *TitleMenuView) ContentSize() (width, height int) {
	return len(logo[0]), len(logo) + 11
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

	optionsX := (screenWidth - 23) / 2
	optionsY := logoY + len(logo) + 2

	defaultGameText := "[SPACE]  Quick-start"
	defaultGameStyle := palette.DefaultGameText
	if v.savedGame != nil {
		defaultGameText = fmt.Sprintf(
			"[SPACE]  Continue [♥ %d] [mines %d]",
			v.savedGame.LivesRemaining(),
			v.savedGame.MinesRemaining(),
		)
		defaultGameStyle = palette.StatusText
	}
	screen.PutStrStyled(optionsX, optionsY, defaultGameText, defaultGameStyle)
	screen.PutStrStyled(optionsX, optionsY+2, "  [1]    H-Expert", palette.ExpertGameText)
	screen.PutStrStyled(optionsX, optionsY+3, "  [2]    H-Big", palette.BigGameText)
	screen.PutStrStyled(optionsX, optionsY+4, "  [3]    Classic Easy", palette.ClassicGameText)
	screen.PutStrStyled(optionsX, optionsY+5, "  [4]    Classic Medium", palette.ClassicGameText)
	screen.PutStrStyled(optionsX, optionsY+6, "  [5]    Classic Expert", palette.ClassicGameText)
	screen.PutStrStyled(optionsX, optionsY+8, " [ESC]   Exit", palette.ExitText)
}

func (v *TitleMenuView) OnInput(key tcell.Key, rune rune) {
	switch key {
	case tcell.KeyEscape:
		v.ui.popView()
	case tcell.KeyEnter:
		v.startDefaultGame()
	default:
		switch rune {
		case ' ':
			v.startDefaultGame()
		case '1':
			v.startGame(newExpertGameFactory())
		case '2':
			v.startGame(v.newBigGameFactory())
		case '3':
			v.startGame(newClassicGameFactory(9, 9, 10))
		case '4':
			v.startGame(newClassicGameFactory(16, 16, 40))
		case '5':
			v.startGame(newClassicGameFactory(30, 16, 99))
		}
	}
}

func (v *TitleMenuView) startDefaultGame() {
	if v.savedGame != nil {
		v.startGame(newExistingGameFactory(v.savedGame))
	} else {
		v.startGame(newExpertGameFactory())
	}
}

func (v *TitleMenuView) startGame(gameFactory GameFactory) {
	v.ui.pushView(newGameView(v.ui, gameFactory, game.DefaultSavePath()))
}
