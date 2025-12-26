package ui

import (
	"fmt"
	"unicode"

	"github.com/gdamore/tcell/v2"
)

type (
	View interface {
		ContentSize() (width, height int)
		Draw(screen tcell.Screen)
		OnInput(key tcell.Key, rune rune)
	}

	Ui struct {
		views  []View
		screen tcell.Screen
	}
)

func NewUi() *Ui {
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}

	err = screen.Init()
	if err != nil {
		panic(err)
	}

	ui := &Ui{
		views:  make([]View, 0),
		screen: screen,
	}
	ui.pushView(newDifficultySelectView(ui))
	return ui
}

func (u *Ui) Loop() {
	for {
		u.refresh()
		switch event := u.screen.PollEvent().(type) {
		case *tcell.EventResize:
			u.fullRefresh()
		case *tcell.EventKey:
			u.topView().OnInput(event.Key(), unicode.ToLower(event.Rune()))
		}
	}
}

func (u *Ui) topView() View {
	return u.views[len(u.views)-1]
}

func (u *Ui) refresh() {
	view := u.topView()

	screenWidth, screenHeight := u.screen.Size()
	contentWidth, contentHeight := view.ContentSize()
	if contentWidth <= screenWidth && contentHeight <= screenHeight {
		view.Draw(u.screen)
	} else {
		message := fmt.Sprintf("Terminal too small (required %dx%d)", contentWidth, contentHeight)
		messageX := (screenWidth - len(message)) / 2
		messageY := screenHeight / 2
		u.screen.PutStrStyled(messageX, messageY, message, defaultPalette.HintText)
	}

	u.screen.Show()
}

func (u *Ui) fullRefresh() {
	u.screen.Fill(' ', defaultPalette.Blank)
	u.refresh()
}

func (u *Ui) pushView(view View) {
	u.views = append(u.views, view)
	u.fullRefresh()
}

func (u *Ui) popView() {
	if len(u.views) > 1 {
		u.views = u.views[:len(u.views)-1]
	}
	u.fullRefresh()
}
