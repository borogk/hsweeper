package game

import (
	"encoding/json"
)

// Snapshot represents a game state, which can be restored and continued from.
type Snapshot struct {
	Status                    Status
	Width                     int
	Height                    int
	MinesToPlant              int
	HeartsToPlant             int
	LivesLeft                 int
	HeartsLeft                int
	MineLocations             []int
	RevealedLocations         []int
	UncollectedHeartLocations []int
	FlaggedLocations          []int
	QuestionedLocations       []int
}

// Encode converts the snapshot into bytes representation.
func (s *Snapshot) Encode() []byte {
	bytes, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return bytes
}

// DecodeSnapshot loads a snapshot from previously encoded bytes data.
func DecodeSnapshot(bytes []byte) (*Snapshot, error) {
	snapshot := &Snapshot{}
	err := json.Unmarshal(bytes, &snapshot)
	return snapshot, err
}
