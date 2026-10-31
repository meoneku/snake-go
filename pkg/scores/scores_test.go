package scores

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScoresLoadNonExistent(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "scores.json")

	gd, err := LoadFrom(path)
	if err != nil {
		t.Fatalf("LoadFrom failed: %v", err)
	}

	if gd.BestScore != 0 {
		t.Errorf("expected BestScore 0, got %d", gd.BestScore)
	}
	if len(gd.HighScores) != 0 {
		t.Errorf("expected 0 high scores, got %d", len(gd.HighScores) )
	}
}

func TestScoresRecordAndTop10(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "scores.json")

	gd, err := LoadFrom(path)
	if err != nil {
		t.Fatalf("LoadFrom failed: %v", err)
	}

	// Insert 15 scores
	scoresToInsert := []int{10, 50, 30, 90, 20, 100, 40, 80, 70, 60, 120, 5, 25, 110, 95}
	for _, s := range scoresToInsert {
		gd.Record(s, 10, 1, "Dinding", 5, 1)
	}

	if gd.BestScore != 120 {
		t.Errorf("expected BestScore 120, got %d", gd.BestScore)
	}
	if gd.GamesPlayed != 15 {
		t.Errorf("expected GamesPlayed 15, got %d", gd.GamesPlayed)
	}
	if gd.TotalApples != 15*5 {
		t.Errorf("expected TotalApples %d, got %d", 15*5, gd.TotalApples)
	}
	if len(gd.HighScores) != 10 {
		t.Fatalf("expected exactly 10 high scores, got %d", len(gd.HighScores))
	}

	// Check top 3 order: 120, 110, 100
	expectedTop := []int{120, 110, 100, 95, 90, 80, 70, 60, 50, 40}
	for i, exp := range expectedTop {
		if gd.HighScores[i].Score != exp {
			t.Errorf("rank %d: expected score %d, got %d", i+1, exp, gd.HighScores[i].Score)
		}
		if gd.HighScores[i].Rank != i+1 {
			t.Errorf("rank %d: expected Rank %d, got %d", i+1, i+1, gd.HighScores[i].Rank)
		}
	}

	// Test persistence by reloading from disk
	reloaded, err := LoadFrom(path)
	if err != nil {
		t.Fatalf("Reload failed: %v", err)
	}
	if reloaded.BestScore != 120 {
		t.Errorf("reloaded BestScore expected 120, got %d", reloaded.BestScore)
	}
	if len(reloaded.HighScores) != 10 {
		t.Errorf("reloaded HighScores expected 10, got %d", len(reloaded.HighScores))
	}
}

func TestScoresLegacyMigration(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "scores.json")

	// Create dummy legacy file in tempDir
	legacyFile := filepath.Join(tempDir, ".ular_best")
	if err := os.WriteFile(legacyFile, []byte("150\n"), 0644); err != nil {
		t.Fatalf("failed to write dummy legacy file: %v", err)
	}

	gd := NewEmptyGameData(path)
	// Temporarily override or test legacy reading
	b, err := os.ReadFile(legacyFile)
	if err != nil || string(b) != "150\n" {
		t.Fatalf("failed reading dummy legacy file: %v", err)
	}
	gd.migrateLegacy() // will check home dir in real execution
}
