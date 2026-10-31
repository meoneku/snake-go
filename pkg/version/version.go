package version

import "fmt"

var (
	// Version holds the current semantic version of the game.
	// Can be overridden at build time via -ldflags.
	Version = "1.0.1"

	// GitCommit holds the short git commit hash.
	GitCommit = "dev"

	// BuildDate holds the build timestamp.
	BuildDate = "unknown"
)

// Info returns formatted version, commit, and build date information.
func Info() string {
	return fmt.Sprintf("Snake Go v%s (%s, built %s)", Version, GitCommit, BuildDate)
}

// Short returns "v" followed by Version (e.g. "v1.0.0").
func Short() string {
	return "v" + Version
}
