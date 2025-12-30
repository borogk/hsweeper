package game

import (
	"os"
	"path"
	"sync"
	"time"
)

// AutoSaver periodically saves the game on disk.
type AutoSaver struct {
	game        *Game
	savePath    string
	ticker      *time.Ticker
	needsToSave bool
	sync.Mutex
}

// DefaultSavePath returns default auto-save path.
func DefaultSavePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = os.TempDir()
	}

	return path.Join(homeDir, ".hsweeper", "autosave.json")
}

// NewAutoSaver creates an auto-saver for specified game and save path.
func NewAutoSaver(game *Game, savePath string) *AutoSaver {
	// Make sure we have a folder to save into
	_ = os.MkdirAll(path.Dir(savePath), 0700)

	s := &AutoSaver{
		game:        game,
		savePath:    savePath,
		ticker:      time.NewTicker(5 * time.Second),
		needsToSave: true,
	}

	// Immediately take over the save file
	s.save()

	// Handle the ticker
	go func() {
		for {
			_ = <-s.ticker.C
			s.Lock()
			s.save()
			s.Unlock()
		}
	}()

	return s
}

// DeferSave remembers that the game needs saving during the next cycle.
func (s *AutoSaver) DeferSave() {
	s.Lock()
	s.needsToSave = true
	s.Unlock()
}

// Finalize stops the autosaving and forcefully saves one last time.
func (s *AutoSaver) Finalize() {
	s.Lock()
	s.ticker.Stop()
	s.save()
	s.Unlock()
}

// Persists the game state on disk. Error handling is disabled to not crash the program and keep playing.
func (s *AutoSaver) save() {
	if s.needsToSave {
		if !s.game.IsFinished() {
			data := s.game.Save().Encode()
			_ = os.WriteFile(s.savePath, data, 0600)
		} else {
			// Delete finished game from disk
			_ = os.Remove(s.savePath)
		}
		s.needsToSave = false
	}
}
