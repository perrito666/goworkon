package actions

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/perrito666/goworkon/environment"
	"github.com/perrito666/goworkon/goinstalls"
	"github.com/perrito666/goworkon/goswitch"
	"github.com/pkg/errors"
)

func configGet(basePath, name string) (environment.Config, error) {
	cfgs, err := environment.LoadConfig(basePath)
	if err != nil {
		return environment.Config{}, errors.WithStack(err)
	}
	env, ok := cfgs[name]
	if !ok {
		return environment.Config{}, errors.Errorf("environment %q not found", name)
	}
	return env, nil
}

// Switch changes the environment to the specified one
// if it exists, otherwise its a noop and returns error.
func Switch(basePath, installName string) error {
	cfgs, err := environment.LoadConfig(basePath)
	if err != nil {
		return errors.WithStack(err)
	}
	env, err := configGet(basePath, installName)
	if err != nil {
		return errors.WithStack(err)
	}

	settings, err := environment.LoadSettings(basePath)
	if err != nil {
		return errors.Wrap(err, "loading settings")
	}

	extraBins := []string{}
	for _, cfg := range cfgs {
		if cfg.GlobalBin {
			extraBins = append(extraBins, filepath.Join(cfg.GoPath, "bin"))
		}
	}
	goswitch.Switch(basePath, env, installName == settings.Default, extraBins)

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

func extractEnvironment(attribute string) (string, string, error) {
	parts := strings.Split(attribute, "@")
	l := len(parts)
	if l == 2 {
		return parts[0], parts[1], nil
	}
	if l == 1 {
		return "", parts[0], nil
	}
	return "", "", errors.Errorf("%q is not a valid attribute", attribute)
}

// Set sets <attribute> to <value> in the correct setting or returns an error.
func Set(attribute, value, baseFolder string) error {
	environmentName, attribute, err := extractEnvironment(attribute)
	if err != nil {
		return errors.Errorf("setting %q to %q", attribute, value)
	}
	if environmentName != "" {
		cfg, err := configGet(baseFolder, environmentName)
		if err != nil {
			return errors.WithStack(err)
		}
		return errors.Wrapf(cfg.Set(attribute, value, baseFolder), "stting %q to %q", attribute, value)

		return errors.New("setting environment attributes is not implemented")
	}

	settings, err := environment.LoadSettings(baseFolder)
	if err != nil {
		return errors.Wrap(err, "loading settings")
	}

	return errors.Wrapf(settings.Set(attribute, value, baseFolder), "setting %q to %q", attribute, value)
}
