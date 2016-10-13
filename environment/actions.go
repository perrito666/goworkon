package environment

import (
	"os"
	"path/filepath"

	"github.com/perrito666/goworkon/goinstalls"
	"github.com/pkg/errors"
)

// Switch changes the environment to the specified one
// if it exists, otherwise its a noop and returns error.
func Switch(installName string) error {
	return nil
}

func ensureVersionInstalled(goVersion, basePath, goroot string) error {
	goFolder := filepath.Join(basePath, goVersion)
	if _, err := os.Stat(goFolder); err == nil {
		return nil
	}
	versions, err := goinstalls.OnlineAvailableVersions()
	if err != nil {
		return errors.WithStack(err)
	}
	for k, v := range versions {
		if k.CommonVersionString() == goVersion {
			goFolder = filepath.Join(basePath, k.String())
			if _, err := os.Stat(goFolder); err == nil {
				return nil
			}
			err = goinstalls.InstallVersion(k, v, basePath, goroot)
			return errors.WithStack(err)
		}
	}
	return errors.Errorf("go version %q not found", goVersion)

}

// Create creates the an environment with the passed name
// in the passed go version, if it exists its a noop and
// returns an error.
func Create(installName, goVersion, basePath string, settings Settings) error {
	return ensureVersionInstalled(goVersion, basePath, settings.Goroot)
}
