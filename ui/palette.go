package ui

import (
	"github.com/borogk/hsweeper/game"
	"github.com/gdamore/tcell/v2"
)

// Palette defines colors for all graphics elements in the game.
// ANSI 256-color palette is used for better terminal compatibility.
type Palette struct {
	Blank                 tcell.Style
	PlainText             tcell.Style
	Logo                  tcell.Style
	LogoSecondary         tcell.Style
	ExpertGameText        tcell.Style
	BigGameText           tcell.Style
	ClassicGameText       tcell.Style
	ExitText              tcell.Style
	Border                tcell.Style
	ReadyText             tcell.Style
	StatusText            tcell.Style
	LoseText              tcell.Style
	WinText               tcell.Style
	Cursor                tcell.Style
	Unrevealed            tcell.Style
	Flag                  tcell.Style
	Question              tcell.Style
	Heart                 tcell.Style
	RevealedMine          tcell.Style
	UnrevealedMine        tcell.Style
	RevealUnrevealedFlash tcell.Style
	RevealFlagFlash       tcell.Style
	BlastFlash            tcell.Style
	Numbers               []tcell.Style
}

var defaultPalette = Palette{
	Blank:                 style(),
	PlainText:             style(255),
	Logo:                  style(232, 160),
	LogoSecondary:         style(232, 26),
	ExpertGameText:        style(226),
	BigGameText:           style(84),
	ClassicGameText:       style(50),
	ExitText:              style(255),
	Border:                style(236),
	ReadyText:             style(255),
	StatusText:            style(196),
	LoseText:              style(196),
	WinText:               style(40),
	Cursor:                style(232, 84),
	Unrevealed:            style(236, 234),
	Flag:                  style(220, 234),
	Question:              style(255, 234),
	Heart:                 style(196),
	RevealedMine:          style(232, 196),
	UnrevealedMine:        style(196),
	RevealUnrevealedFlash: style(226, 234),
	RevealFlagFlash:       style(231, 234),
	BlastFlash:            style(232, 196),
	Numbers: []tcell.Style{
		style(232),
		style(33),
		style(84),
		style(196),
		style(213),
		style(88),
		style(27),
		style(92),
		style(244),
	},
}

var lostPalette = defaultPalette
var wonPalette = defaultPalette

func init() {
	lostPalette.Border = style(52)
	lostPalette.Flag = style(196)
	lostPalette.Question = style(196)
	lostPalette.Unrevealed = style(52)
	lostPalette.UnrevealedMine = style(196)
	lostPalette.Numbers = []tcell.Style{style(52)}

	wonPalette.Border = style(22)
	wonPalette.Flag = style(84)
	wonPalette.Question = style(84)
	wonPalette.UnrevealedMine = style(84)
	wonPalette.Numbers = []tcell.Style{style(22)}
	wonPalette.RevealUnrevealedFlash = style(22)
	wonPalette.RevealFlagFlash = style(84)
	wonPalette.BlastFlash = style(232, 84)
}

// Returns palette fitting for current game status.
func gamePalette(g *game.Game) Palette {
	if g.Status() == game.StatusLost {
		return lostPalette
	} else if g.Status() == game.StatusWon {
		return wonPalette
	}

	return defaultPalette
}

func style(colors ...byte) tcell.Style {
	foreground := byte(232)
	background := byte(232)
	if len(colors) > 0 {
		foreground = colors[0]
	}
	if len(colors) > 1 {
		background = colors[1]
	}

	return tcell.StyleDefault.Foreground(ansiColor(foreground)).Background(ansiColor(background))
}

func ansiColor(color byte) tcell.Color {
	return tcell.ColorValid | tcell.Color(color)
}
