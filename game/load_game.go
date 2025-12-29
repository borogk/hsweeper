package game

import "os"

// LoadGame loads gave from disk, may return nil if it couldn't load.
func LoadGame(path string) *Game {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	snapshot, err := DecodeSnapshot(data)
	if err != nil {
		return nil
	}

	return RestoreGame(snapshot)
}
