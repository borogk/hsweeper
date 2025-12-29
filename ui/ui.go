package ui

import (
	"fmt"
	"os"
	"unicode"

	"github.com/gdamore/tcell/v2"
)

type (
	// View implementations are responsible for (almost) the entire screen input and graphics.
	View interface {
		OnActivate()
		OnDeactivate()
		ContentSize() (width, height int)
		Draw(screen tcell.Screen)
		OnInput(key tcell.Key, rune rune)
	}

	// Ui encapsulates all game graphics and input.
	Ui struct {
		views  []View
		screen tcell.Screen
	}
)

// NewUiWithTitleMenu creates new UI with title menu as its starting view.
func NewUiWithTitleMenu() *Ui {
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
	ui.pushView(newTitleMenuView(ui))
	return ui
}

// Loop processes all input and graphics in a loop.
func (u *Ui) Loop() {
	for {
		u.refresh()
		switch event := u.screen.PollEvent().(type) {
		case *tcell.EventResize:
			u.fullRefresh()
		case *tcell.EventKey:
			if event.Key() == tcell.KeyCtrlC {
				u.exit()
			} else {
				u.topView().OnInput(event.Key(), unicode.ToLower(event.Rune()))
			}
		}
	}
}

// Returns current top view in the stack.
func (u *Ui) topView() View {
	return u.views[len(u.views)-1]
}

// Refreshes the graphics.
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
		u.screen.PutStrStyled(messageX, messageY, message, defaultPalette.PlainText)
	}

	u.screen.Show()
}

// Refreshes the graphics, making sure the screen is fully wiped first.
func (u *Ui) fullRefresh() {
	u.screen.Fill(' ', defaultPalette.Blank)
	u.refresh()
}

// Puts a new view as the current.
func (u *Ui) pushView(view View) {
	u.views = append(u.views, view)
	view.OnActivate()
	u.fullRefresh()
}

// Removes current top view, replacing it with the view underneath. If no more views left - exit the program.
func (u *Ui) popView() {
	u.topView().OnDeactivate()
	u.views = u.views[:len(u.views)-1]
	if len(u.views) > 0 {
		u.topView().OnActivate()
		u.fullRefresh()
	} else {
		u.exit()
	}
}

// Gracefully exits the program.
func (u *Ui) exit() {
	for _, view := range u.views {
		view.OnDeactivate()
	}

	u.screen.Fini()
	os.Exit(0)
}
