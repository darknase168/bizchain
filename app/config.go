package app

import (
	"os"
	"path/filepath"
)

const (
	// DefaultNodeHome sets the default home directory name for bizchaind
	DefaultNodeHome = ".bizchain"
)

// GetDefaultNodeHome returns the default node home directory path
func GetDefaultNodeHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return DefaultNodeHome
	}
	return filepath.Join(home, DefaultNodeHome)
}
