package paths

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

const (
	// This is very not cross-os it only works for unixes, help is
	// welcome for windowses.

	// UNIXHOMEVAR holds the name of the home variable in unix.
	UNIXHOMEVAR = "HOME"
	// UNIXDEFAULTXDGDATAREL holds the default relative path for
	// xdg data (relative to $HOME)
	UNIXDEFAULTXDGDATAREL = ".local/share"
	// UNIXXDGDATAHOME holds the name of the variable that might contain
	// the path to XDG data home.
	UNIXXDGDATAHOME = "XDG_DATA_HOME"

	// PATHSEPARATOR holds the character used to separate PATH members.
	PATHSEPARATOR = ":"

	// GOWORKONNAME holds the name of goworkon xdg folder
	GOWORKONNAME = "goworkon"
	// CONFIGSFOLDER holds the name of the envs configuration folder
	// inside goworkon xdg home.
	CONFIGSFOLDER = "configs"
	// INSTALLSFOLDER holds the name of the go install s folder insde
	// goworkon xdg home.
	INSTALLSFOLDER = "installs"
)

// XdgData returns the most likely place for XDG data to be
// this will most likely notwork in windows.
func XdgData() (string, error) {
	xdgConfig := os.Getenv(UNIXXDGDATAHOME)
	if xdgConfig != "" {
		return filepath.Join(xdgConfig, GOWORKONNAME), nil
	}
	home := os.Getenv(UNIXHOMEVAR)
	if home == "" {
		return "", errors.Errorf("cannot determine $%s", UNIXHOMEVAR)
	}
	// XDG standard says $HOME/.local/share/$PROJECTNAME is the
	// fallback if the variable is not set.
	return filepath.Join(home, UNIXDEFAULTXDGDATAREL, GOWORKONNAME), nil
}

// XdgDataConfig returns the folder where config should be stored for
// environments.
func XdgDataConfig() (string, error) {
	xdgDataDir, err := XdgData()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return filepath.Join(xdgDataDir, CONFIGSFOLDER), nil
}

// XdgDataGoInstalls returns the folder where go installs should live.
func XdgDataGoInstalls() (string, error) {
	xdgDataDir, err := XdgData()
	if err != nil {
		return "", errors.WithStack(err)
	}
	return filepath.Join(xdgDataDir, INSTALLSFOLDER), nil
}

// XdgDataGoInstallsBinForVerson returns the bin path of the given go version.
func XdgDataGoInstallsBinForVerson(goVersion string) (string, error) {
	installs, err := XdgDataGoInstalls()
	if err != nil {
		return "", errors.Wrapf(err, "trying to determine go bin path for %q", goVersion)
	}
	return filepath.Join(installs, goVersion, "go", "bin"), nil
}

// GoPathBin is a convenience function that returns the bin folder of  a
// given GOPATH.
func GoPathBin(goPath string) string {
	return filepath.Join(goPath, "bin")
}

// PATHInsert returns a string representing PATH after inserting
// <newMembers> at the beginning and removing them from other positions
// if they exist.
func PATHInsert(currentPath string, newMembers ...string) string {
	pathMap := make(map[string]bool, len(newMembers))
	for _, member := range newMembers {
		pathMap[member] = true
	}
	pathMembers := strings.Split(currentPath, PATHSEPARATOR)
	for _, oldPathMember := range pathMembers {
		_, ok := pathMap[oldPathMember]
		if !ok {
			newMembers = append(newMembers, oldPathMember)
		}
	}
	return strings.Join(newMembers, PATHSEPARATOR)
}
