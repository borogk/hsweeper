package ui

import (
	"fmt"
	"minesweeper/game"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

type (
	Effect struct {
		x       int
		y       int
		style   tcell.Style
		expired bool
	}

	GameView struct {
		ui           *Ui
		game         *game.Game
		cx           int
		cy           int
		effects      []*Effect
		effectsMutex sync.Mutex
	}
)

func newGameView(ui *Ui, game *game.Game) *GameView {
	return &GameView{
		ui:      ui,
		game:    game,
		cx:      game.Width() / 2,
		cy:      game.Height() / 2,
		effects: make([]*Effect, 0),
	}
}

func (v *GameView) ContentSize() (width, height int) {
	return v.game.Width()*3 + 2, v.game.Height() + 4
}

func (v *GameView) Draw(screen tcell.Screen) {
	screenWidth, screenHeight := screen.Size()
	contentWidth, contentHeight := v.ContentSize()
	palette := gamePalette(v.game)

	offsetX := (screenWidth - contentWidth) / 2
	offsetY := (screenHeight - contentHeight) / 2

	statusMessage, statusStyle, statusCentered := v.statusAppearance(palette)
	statusX := offsetX + 1
	statusY := offsetY + 1
	if statusCentered {
		statusX += (v.game.Width()*3 - len(statusMessage)) / 2
	}
	screen.PutStrStyled(0, statusY, fmt.Sprintf("%*s", screenWidth, ""), palette.Blank)
	screen.PutStrStyled(statusX, statusY, statusMessage, statusStyle)

	borderLeft := offsetX
	borderRight := borderLeft + v.game.Width()*3 + 1
	borderTop := statusY + 1
	borderBottom := borderTop + v.game.Height() + 1
	screen.Put(borderLeft, borderTop, "┌", palette.Border)
	screen.Put(borderRight, borderTop, "┐", palette.Border)
	screen.Put(borderLeft, borderBottom, "└", palette.Border)
	screen.Put(borderRight, borderBottom, "┘", palette.Border)
	for x := borderLeft + 1; x < borderRight; x++ {
		screen.Put(x, borderTop, "─", palette.Border)
		screen.Put(x, borderBottom, "─", palette.Border)
	}
	for y := borderTop + 1; y < borderBottom; y++ {
		screen.Put(borderLeft, y, "│", palette.Border)
		screen.Put(borderRight, y, "│", palette.Border)
	}

	printCell := func(x, y int, symbol string, style tcell.Style) {
		cellX := borderLeft + 1 + x*3
		cellY := borderTop + 1 + y
		screen.PutStrStyled(cellX, cellY, symbol, style)
	}

	for x := 0; x < v.game.Width(); x++ {
		for y := 0; y < v.game.Height(); y++ {
			symbol, style := v.cellAppearance(x, y, palette)
			printCell(x, y, symbol, style)
		}
	}

	v.effectsMutex.Lock()
	validEffects := make([]*Effect, 0)
	for _, effect := range v.effects {
		if !effect.expired {
			symbol, _ := v.cellAppearance(effect.x, effect.y, palette)
			printCell(effect.x, effect.y, symbol, effect.style)
			validEffects = append(validEffects, effect)
		}
	}
	v.effects = validEffects
	v.effectsMutex.Unlock()
}

func (v *GameView) OnInput(key tcell.Key, rune rune) {
	switch key {
	case tcell.KeyLeft:
		v.moveCursor(-1, 0)
	case tcell.KeyRight:
		v.moveCursor(1, 0)
	case tcell.KeyUp:
		v.moveCursor(0, -1)
	case tcell.KeyDown:
		v.moveCursor(0, 1)
	case tcell.KeyEnter:
		v.actionButton()
	case tcell.KeyEscape:
		v.ui.popView()
	case tcell.KeyDelete, tcell.KeyBackspace:
		v.game.ClearFlagAndQuestion(v.cx, v.cy)
	default:
		switch rune {
		case ' ':
			v.actionButton()
		case 'r':
			revealResult := v.game.Reveal(v.cx, v.cy)
			if revealResult == game.RevealResultBlast {
				v.startBlastFlashEffect()
			}
		case 'f':
			v.game.ToggleFlag(v.cx, v.cy)
		case 'q', '?':
			v.game.ToggleQuestion(v.cx, v.cy)
		}
	}
}

func (v *GameView) statusAppearance(palette Palette) (message string, style tcell.Style, centered bool) {
	switch v.game.Status() {
	case game.StatusReady:
		return "Ready?", palette.ReadyText, true
	case game.StatusStarted:
		const maxLives = 6
		livesString := ""
		livesPadding := maxLives*2 + 3
		for i := 0; i < v.game.LivesRemaining() && i < maxLives; i++ {
			livesString += "♥ "
		}
		if v.game.LivesRemaining() > maxLives {
			livesString += fmt.Sprintf("+%-2d", v.game.LivesRemaining()-maxLives)
		}

		return fmt.Sprintf(
			"%-*s%-*smines:%4d",
			livesPadding,
			livesString,
			v.game.Width()*3-livesPadding-10,
			"",
			v.game.MinesRemaining(),
		), palette.StatusText, false
	case game.StatusLost:
		return "GAME OVER", palette.LoseText, true
	case game.StatusWon:
		return "Well done!", palette.WinText, true
	default:
		panic("Unknown game status")
	}
}

func (v *GameView) cellAppearance(x, y int, palette Palette) (symbol string, style tcell.Style) {
	symbol = "   "
	style = palette.Blank

	cell := v.game.Cell(x, y)
	if cell.IsRevealed() {
		if cell.IsMine() {
			symbol = " * "
			style = palette.RevealedMine
		} else if cell.AdjacentMines() > 0 {
			symbol = fmt.Sprintf(" %c ", rune('0'+cell.AdjacentMines()))
			style = palette.Numbers[cell.AdjacentMines()%len(palette.Numbers)]
		} else if cell.IsHeart() {
			symbol = " ♥ "
			style = palette.Heart
		}
	} else if cell.IsQuestioned() {
		symbol = " ? "
		style = palette.Question
	} else if cell.IsFlagged() {
		symbol = " ⚑ "
		style = palette.Flag
	} else if cell.IsMine() && v.game.IsFinished() {
		symbol = " * "
		style = palette.UnrevealedMine
	} else {
		symbol = " ■ "
		style = palette.Unrevealed
	}

	if !v.game.IsFinished() && x == v.cx && y == v.cy {
		style = palette.Cursor
	}

	return
}

func (v *GameView) moveCursor(dx, dy int) {
	if !v.game.IsFinished() && !v.game.IsOutOfBounds(v.cx+dx, v.cy+dy) {
		v.cx = v.cx + dx
		v.cy = v.cy + dy
	}
}

func (v *GameView) actionButton() {
	cell := v.game.Cell(v.cx, v.cy)
	if v.game.Status() == game.StatusReady {
		v.game.Reveal(v.cx, v.cy)

	} else if v.game.Status() == game.StatusStarted {
		if !cell.IsRevealed() {
			v.game.ToggleFlag(v.cx, v.cy)
		} else if cell.AdjacentMines() > 0 {
			revealResult := v.game.RevealAdjacent(v.cx, v.cy)
			if revealResult != game.RevealResultBlast {
				v.startRevealFlashEffect()
			} else {
				v.startBlastFlashEffect()
			}
		} else if cell.IsHeart() {
			v.game.Pickup(v.cx, v.cy)
		}
	} else {
		v.ui.popView()
	}
}

func (v *GameView) startRevealFlashEffect() {
	palette := gamePalette(v.game)
	duration := 300 * time.Millisecond

	unrevealedEffects := v.innerFlashEffects(palette.RevealUnrevealedFlash)
	unrevealedFilter := func(x, y int) bool {
		cell := v.game.Cell(x, y)
		return !cell.IsRevealed() && !cell.IsFlagged() && !cell.IsQuestioned()
	}
	v.startEffects(duration, unrevealedEffects, unrevealedFilter)

	flaggedEffects := v.innerFlashEffects(palette.RevealFlagFlash)
	flaggedFilter := func(x, y int) bool { return v.game.Cell(x, y).IsFlagged() }
	v.startEffects(duration, flaggedEffects, flaggedFilter)
}

func (v *GameView) startBlastFlashEffect() {
	palette := gamePalette(v.game)

	innerDuration := 300 * time.Millisecond
	outerDuration := 200 * time.Millisecond
	innerEffects := v.innerFlashEffects(palette.BlastFlash)
	outerEffects := v.outerFlashEffects(palette.BlastFlash)
	filter := func(x, y int) bool { return true }

	v.startEffects(innerDuration, innerEffects, filter)
	v.startEffects(outerDuration, outerEffects, filter)
}

func (v *GameView) innerFlashEffects(style tcell.Style) []*Effect {
	return []*Effect{
		{x: v.cx - 1, y: v.cy - 1, style: style},
		{x: v.cx - 1, y: v.cy, style: style},
		{x: v.cx - 1, y: v.cy + 1, style: style},
		{x: v.cx, y: v.cy - 1, style: style},
		{x: v.cx, y: v.cy, style: style},
		{x: v.cx, y: v.cy + 1, style: style},
		{x: v.cx + 1, y: v.cy - 1, style: style},
		{x: v.cx + 1, y: v.cy, style: style},
		{x: v.cx + 1, y: v.cy + 1, style: style},
	}
}

func (v *GameView) outerFlashEffects(style tcell.Style) []*Effect {
	return []*Effect{
		{x: v.cx - 2, y: v.cy, style: style},
		{x: v.cx + 2, y: v.cy, style: style},
		{x: v.cx, y: v.cy - 2, style: style},
		{x: v.cx, y: v.cy + 2, style: style},
	}
}

func (v *GameView) startEffects(expireAfter time.Duration, effects []*Effect, filter func(int, int) bool) {
	v.effectsMutex.Lock()
	for _, effect := range effects {
		if filter(effect.x, effect.y) && !v.game.IsOutOfBounds(effect.x, effect.y) {
			v.effects = append(v.effects, effect)
		}
	}
	v.effectsMutex.Unlock()

	time.AfterFunc(expireAfter, func() {
		v.effectsMutex.Lock()
		for _, effect := range effects {
			effect.expired = true
		}
		v.effectsMutex.Unlock()
		v.ui.refresh()
	})
}
