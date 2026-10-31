package scores

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"banguncandi/pkg/version"
)

const CurrentSchemaVersion = 1

// ScoreEntry represents a single completed game run in the leaderboard.
type ScoreEntry struct {
	Rank      int       `json:"rank"`
	Score     int       `json:"score"`
	Player    string    `json:"player"`
	Length    int       `json:"length"`
	Level     int       `json:"level"`
	Mode      string    `json:"mode"` // "Dinding" or "Tembus"
	Timestamp time.Time `json:"timestamp"`
}

// GameData holds persistent game statistics and the top leaderboard scores.
type GameData struct {
	SchemaVersion int          `json:"schema_version"`
	AppVersion    string       `json:"app_version"`
	BestScore     int          `json:"best_score"`
	GamesPlayed   int          `json:"games_played"`
	TotalApples   int          `json:"total_apples"`
	TotalStars    int          `json:"total_stars"`
	HighScores    []ScoreEntry `json:"high_scores"`

	filePath string     `json:"-"`
	mu       sync.Mutex `json:"-"`
}

// DefaultScoresPath returns the OS-standard path for scores.json.
// Windows: %APPDATA%\SnakeGo\scores.json
// Linux: ~/.config/snake-go/scores.json
// macOS: ~/Library/Application Support/SnakeGo/scores.json
func DefaultScoresPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil || configDir == "" {
		home, err := os.UserHomeDir()
		if err == nil && home != "" {
			return filepath.Join(home, ".snake-go", "scores.json")
		}
		return "scores.json"
	}
	return filepath.Join(configDir, "SnakeGo", "scores.json")
}

// LegacyBestPath returns the path to the legacy ~/.ular_best file.
func LegacyBestPath() string {
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".ular_best")
	}
	return ".ular_best"
}

// NewEmptyGameData creates a default initialized GameData.
func NewEmptyGameData(path string) *GameData {
	return &GameData{
		SchemaVersion: CurrentSchemaVersion,
		AppVersion:    version.Version,
		BestScore:     0,
		GamesPlayed:   0,
		TotalApples:   0,
		TotalStars:    0,
		HighScores:    make([]ScoreEntry, 0, 10),
		filePath:      path,
	}
}

// Load loads the scores from the default OS path, migrating legacy data if necessary.
func Load() (*GameData, error) {
	return LoadFrom(DefaultScoresPath())
}

// LoadFrom loads scores from a specific file path.
func LoadFrom(path string) (*GameData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			gd := NewEmptyGameData(path)
			gd.migrateLegacy()
			_ = gd.Save()
			return gd, nil
		}
		return NewEmptyGameData(path), err
	}

	var gd GameData
	if err := json.Unmarshal(data, &gd); err != nil {
		// Corrupted file, fallback to empty
		fallback := NewEmptyGameData(path)
		return fallback, fmt.Errorf("unmarshal scores.json: %w", err)
	}

	gd.filePath = path
	if gd.HighScores == nil {
		gd.HighScores = make([]ScoreEntry, 0, 10)
	}

	return &gd, nil
}

// migrateLegacy reads ~/.ular_best and imports it if valid.
func (gd *GameData) migrateLegacy() {
	legacyFile := LegacyBestPath()
	b, err := os.ReadFile(legacyFile)
	if err != nil {
		return
	}
	val, err := strconv.Atoi(strings.TrimSpace(string(b)))
	if err != nil || val <= 0 {
		return
	}

	gd.BestScore = val
	gd.GamesPlayed = 1
	gd.HighScores = append(gd.HighScores, ScoreEntry{
		Rank:      1,
		Score:     val,
		Player:    "Pemain",
		Length:    val/10 + 3,
		Level:     1,
		Mode:      "Dinding",
		Timestamp: time.Now(),
	})
}

// Record adds a new game score, updates stats, sorts Top 10, and saves to disk.
// Returns isNewBest (true if score > old best) and rank (1-10 if in leaderboard, 0 otherwise).
func (gd *GameData) Record(score, length, level int, mode string, apples, stars int) (bool, int) {
	gd.mu.Lock()
	defer gd.mu.Unlock()

	gd.GamesPlayed++
	gd.TotalApples += apples
	gd.TotalStars += stars

	isNewBest := false
	if score > gd.BestScore {
		gd.BestScore = score
		isNewBest = true
	}

	newEntry := ScoreEntry{
		Score:     score,
		Player:    "Pemain",
		Length:    length,
		Level:     level,
		Mode:      mode,
		Timestamp: time.Now(),
	}

	gd.HighScores = append(gd.HighScores, newEntry)

	// Sort descending by score, ascending by timestamp
	sort.SliceStable(gd.HighScores, func(i, j int) bool {
		if gd.HighScores[i].Score != gd.HighScores[j].Score {
			return gd.HighScores[i].Score > gd.HighScores[j].Score
		}
		return gd.HighScores[i].Timestamp.Before(gd.HighScores[j].Timestamp)
	})

	// Keep only top 10
	if len(gd.HighScores) > 10 {
		gd.HighScores = gd.HighScores[:10]
	}

	// Re-assign ranks and find newEntry's rank
	rank := 0
	for i := range gd.HighScores {
		gd.HighScores[i].Rank = i + 1
		if gd.HighScores[i].Timestamp.Equal(newEntry.Timestamp) && gd.HighScores[i].Score == newEntry.Score {
			rank = i + 1
		}
	}

	gd.AppVersion = version.Version
	_ = gd.saveLocked()

	return isNewBest, rank
}

// Save writes data to disk atomically.
func (gd *GameData) Save() error {
	gd.mu.Lock()
	defer gd.mu.Unlock()
	return gd.saveLocked()
}

func (gd *GameData) saveLocked() error {
	path := gd.filePath
	if path == "" {
		path = DefaultScoresPath()
		gd.filePath = path
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create scores dir %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(gd, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal scores: %w", err)
	}

	tmpPath := path + ".tmp"
	f, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open tmp file %s: %w", tmpPath, err)
	}

	if _, err := f.Write(data); err != nil {
		f.Close()
		_ = os.Remove(tmpPath)
		return fmt.Errorf("write tmp file: %w", err)
	}

	if err := f.Sync(); err != nil {
		f.Close()
		_ = os.Remove(tmpPath)
		return fmt.Errorf("sync tmp file: %w", err)
	}

	if err := f.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("close tmp file: %w", err)
	}

	// On Windows, target must be removed first if it exists
	_ = os.Remove(path)

	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("atomic rename from %s to %s: %w", tmpPath, path, err)
	}

	return nil
}
