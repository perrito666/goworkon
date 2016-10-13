package actions

import (
	"os"
	"path/filepath"

	"github.com/perrito666/goworkon/environment"
	"github.com/perrito666/goworkon/goinstalls"
	"github.com/perrito666/goworkon/goswitch"
	"github.com/pkg/errors"
)

// Switch changes the environment to the specified one
// if it exists, otherwise its a noop and returns error.
func Switch(basePath, installName string) error {
	cfgs, err := environment.LoadConfig(basePath)
	if err != nil {
		return errors.WithStack(err)
	}
	env, ok := cfgs[installName]
	if !ok {
		return errors.Errorf("environment %q not found", installName)
	}
	goswitch.Switch(basePath, env)

	return nil
}

// Reset will try toreturn the env to its original state.
func Reset() error {
	return goswitch.Reset()
}

func ensureVersionInstalled(goVersion, basePath, goroot string) (goinstalls.Version, error) {
	installFolder := filepath.Join(basePath, "installs")
	goFolder := filepath.Join(installFolder, goVersion)
	if _, err := os.Stat(goFolder); err == nil {
		return goinstalls.VersionFromString(goVersion)
	}
	versions, err := goinstalls.OnlineAvailableVersions()
	if err != nil {
		return goinstalls.Version{}, errors.WithStack(err)
	}
	for k, v := range versions {
		if k.CommonVersionString() == goVersion {
			goFolder = filepath.Join(installFolder, k.String())
			if _, err := os.Stat(goFolder); err == nil {
				return k, nil
			}
			err = goinstalls.InstallVersion(k, v, installFolder, goroot)
			return k, errors.WithStack(err)
		}
	}
	return goinstalls.Version{}, errors.Errorf("go version %q not found", goVersion)

}

// Create creates the an environment with the passed name
// in the passed go version, if it exists its a noop and
// returns an error.
func Create(installName, goVersion, basePath, goPath string, settings environment.Settings) error {
	cfgs, err := environment.LoadConfig(basePath)
	if err != nil {
		return errors.WithStack(err)
	}
	if _, ok := cfgs[installName]; ok {
		return nil
	}
	v, err := ensureVersionInstalled(goVersion, basePath, settings.Goroot)
	if err != nil {
		return errors.WithStack(err)
	}
	c := environment.Config{
		Name:      installName,
		GoVersion: v.String(),
		GoPath:    goPath,
	}
	if err := c.Save(basePath); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
