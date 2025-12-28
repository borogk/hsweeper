package game

type RevealResult byte

const (
	RevealResultBlocked RevealResult = iota
	RevealResultRevealed
	RevealResultBlast
)
