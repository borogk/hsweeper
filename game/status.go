package game

type Status byte

const (
	StatusReady Status = iota
	StatusStarted
	StatusLost
	StatusWon
)
