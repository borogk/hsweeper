package ui

import (
	"fmt"

	"github.com/borogk/hsweeper/game"
	"github.com/gdamore/tcell/v2"
)

type (
	// TitleMenuItem represents a menu item with its action and appearance.
	TitleMenuItem struct {
		action func()
		text   string
		style  tcell.Style
		margin int
	}

	// TitleMenuView is responsible for title menu input and graphics.
	TitleMenuView struct {
		ui        *Ui
		savedGame *game.Game
		items     []TitleMenuItem
		cursor    int
	}
)

// The symbols from this logo are not put directly into terminal, but interpreted as colored text cell backgrounds.
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
	// Preload the auto-save each time menu is activated
	v.savedGame = game.LoadGame(game.DefaultSavePath())
	v.refreshMenuItems()
}

func (v *TitleMenuView) OnDeactivate() {

}

func (v *TitleMenuView) OnInput(key tcell.Key, rune rune) {
	switch key {
	case tcell.KeyDown:
		v.cursor = (v.cursor + 1) % len(v.items)
	case tcell.KeyUp:
		v.cursor = (v.cursor - 1 + len(v.items)) % len(v.items)
	case tcell.KeyEscape:
		v.ui.popView()
	case tcell.KeyEnter:
		v.selectMenuItem()
	default:
		switch rune {
		case ' ':
			v.selectMenuItem()
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

func (v *TitleMenuView) ContentSize() (width, height int) {
	height = len(logo) + 2
	for _, item := range v.items {
		height += 1 + item.margin
	}

	return len(logo[0]), height
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

	itemsX := (screenWidth - 23) / 2
	itemsY := logoY + len(logo) + 2
	for i, item := range v.items {
		screen.PutStrStyled(itemsX, itemsY, item.text, item.style)
		if i == v.cursor {
			screen.PutStrStyled(itemsX-2, itemsY, "▶", item.style)
		} else {
			screen.PutStrStyled(itemsX-2, itemsY, " ", item.style)
		}
		itemsY += 1 + item.margin
	}
}

func (v *TitleMenuView) refreshMenuItems() {
	v.items = make([]TitleMenuItem, 0, 7)

	if v.savedGame != nil {
		v.items = append(v.items, TitleMenuItem{
			text: fmt.Sprintf(
				"Continue [♥ %d] [mines %d]",
				v.savedGame.LivesRemaining(),
				v.savedGame.MinesRemaining(),
			),
			style:  defaultPalette.StatusText,
			action: func() { v.startGame(newExistingGameFactory(v.savedGame)) },
			margin: 1,
		})
	}

	v.items = append(v.items, TitleMenuItem{
		text:   " 1   H-Expert",
		style:  defaultPalette.ExpertGameText,
		action: func() { v.startGame(newExpertGameFactory()) },
	})
	v.items = append(v.items, TitleMenuItem{
		text:   " 2   H-Big",
		style:  defaultPalette.BigGameText,
		action: func() { v.startGame(v.newBigGameFactory()) },
	})
	v.items = append(v.items, TitleMenuItem{
		text:   " 3   Classic Easy",
		style:  defaultPalette.ClassicGameText,
		action: func() { v.startGame(newClassicGameFactory(9, 9, 10)) },
	})
	v.items = append(v.items, TitleMenuItem{
		text:   " 4   Classic Medium",
		style:  defaultPalette.ClassicGameText,
		action: func() { v.startGame(newClassicGameFactory(16, 16, 40)) },
	})
	v.items = append(v.items, TitleMenuItem{
		text:   " 5   Classic Expert",
		style:  defaultPalette.ClassicGameText,
		action: func() { v.startGame(newClassicGameFactory(30, 16, 99)) },
		margin: 1,
	})
	v.items = append(v.items, TitleMenuItem{
		text:   "ESC  Exit",
		style:  defaultPalette.ExitText,
		action: func() { v.ui.popView() },
	})

	v.cursor = 0
}

func (v *TitleMenuView) selectMenuItem() {
	v.items[v.cursor].action()
}

func (v *TitleMenuView) startGame(gameFactory GameFactory) {
	v.ui.pushView(newGameView(v.ui, gameFactory, game.DefaultSavePath()))
}
