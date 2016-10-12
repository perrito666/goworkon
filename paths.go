package main

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// xdgData returns the most likely place for XDG data to be
// this will most likely notwork in windows.
func xdgData() (string, error) {
	xdgConfig := os.Getenv("XDG_DATA_HOME")
	if xdgConfig != "" {
		return filepath.Join(xdgConfig, "goworkon"), nil
	}
	home := os.Getenv("HOME")
	if home == "" {
		return "", errors.New("cannot determine $HOME")
	}
	// XDG standard says $HOME/.local/share/$PROJECTNAME is the
	// fallback if the variable is not set.
	return filepath.Join(home, ".local", "share", "goworkon"), nil
}

func xdgDataConfig() (string, error) {
	xdgDataDir, err := xdgData()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return filepath.Join(xdgDataDir, "configs"), nil
}

func xdgDataGoInstalls() (string, error) {
	xdgDataDir, err := xdgData()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return filepath.Join(xdgDataDir, "installs"), nil
}
