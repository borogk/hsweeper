package game

import (
	"testing"
)

func TestSnapshot_EncodeDecode(t *testing.T) {
	snapshot := &Snapshot{
		Status:                    StatusStarted,
		Width:                     50,
		Height:                    40,
		MinesToPlant:              10,
		HeartsToPlant:             4,
		LivesLeft:                 3,
		HeartsLeft:                2,
		MineLocations:             []int{4, 5, 6, 7},
		RevealedLocations:         []int{1, 2, 3},
		UncollectedHeartLocations: []int{8, 9},
		FlaggedLocations:          []int{4, 5, 6},
		QuestionedLocations:       []int{7},
	}

	bytes := snapshot.Encode()
	decodedSnapshot, _ := DecodeSnapshot(bytes)

	assertEquals(t, decodedSnapshot, snapshot)
}
