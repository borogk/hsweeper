package ui

import (
	"hsweeper/game"

	"github.com/gdamore/tcell/v2"
)

type Palette struct {
	Blank                 tcell.Style
	Logo                  tcell.Style
	LogoSecondary         tcell.Style
	HintText              tcell.Style
	ExpertGameText        tcell.Style
	BigGameText           tcell.Style
	ClassicGameText       tcell.Style
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

const blankColor = 0x0C0C0C

func rgb(color uint64) tcell.Color {
	return tcell.ColorIsRGB | tcell.ColorValid | tcell.Color(color)
}

func style(foreground uint64) tcell.Style {
	return styleWithBackground(foreground, blankColor)
}

func styleWithBackground(foreground, background uint64) tcell.Style {
	return tcell.StyleDefault.Foreground(rgb(foreground)).Background(rgb(background))
}

var defaultPalette = Palette{
	Blank:                 style(blankColor),
	Logo:                  styleWithBackground(blankColor, 0xCC1C45),
	LogoSecondary:         styleWithBackground(blankColor, 0x274FBC),
	HintText:              style(0xB3B3B3),
	ExpertGameText:        style(0xFFA000),
	BigGameText:           style(0xFF0000),
	ClassicGameText:       style(0xD0D0FF),
	Border:                style(0x352D66),
	ReadyText:             style(0xD3D3D3),
	StatusText:            style(0xFF0000),
	LoseText:              style(0xFF0000),
	WinText:               style(0x90EE90),
	Cursor:                styleWithBackground(0x000000, 0x9ACD32),
	Unrevealed:            styleWithBackground(0x352D66, 0x201B3D),
	Flag:                  style(0xFFA500),
	Question:              style(0xFFFFFF),
	Heart:                 style(0xFF0000),
	RevealedMine:          styleWithBackground(0x000000, 0xFF0000),
	UnrevealedMine:        style(0xFF0000),
	RevealUnrevealedFlash: styleWithBackground(0xFFFF00, 0x201B3D),
	RevealFlagFlash:       style(0xFFFF00),
	BlastFlash:            styleWithBackground(0x000000, 0xFF0000),
	Numbers: []tcell.Style{
		style(0x000000),
		style(0x0080FF),
		style(0x90EE90),
		style(0xFF0000),
		style(0xEE82EE),
		style(0xA52A2A),
		style(0x008B8B),
		style(0x8A2BE2),
		style(0x808080),
	},
}

func gamePalette(g *game.Game) Palette {
	palette := defaultPalette
	if g.Status() == game.StatusLost {
		palette.Border = style(0x480000)
		palette.Flag = style(0xFF0000)
		palette.Question = style(0xFF0000)
		palette.Unrevealed = style(0x480000)
		palette.UnrevealedMine = style(0xFF0000)
		palette.Numbers = []tcell.Style{style(0x480000)}
	} else if g.Status() == game.StatusWon {
		palette.Border = style(0x006400)
		palette.Flag = style(0x90EE90)
		palette.Question = style(0x90EE90)
		palette.UnrevealedMine = style(0x90EE90)
		palette.Numbers = []tcell.Style{style(0x006400)}
	}

	return palette
}
