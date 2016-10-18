package actions

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/perrito666/goworkon/environment"
	"github.com/perrito666/goworkon/goinstalls"
	"github.com/perrito666/goworkon/goswitch"
	"github.com/perrito666/goworkon/paths"
	"github.com/pkg/errors"
)

func globalBins() ([]string, error) {
	basePath, err := paths.XdgDataConfig()
	if err != nil {
		return nil, errors.Wrap(err, "retrieving config files")
	}
	cfgs, err := environment.LoadConfig(basePath)
	if err != nil {
		return nil, errors.Wrap(err, "loading configs")
	}
	globalBin := []string{}
	for _, cfg := range cfgs {
		if cfg.GlobalBin {
			globalBin = append(globalBin, paths.GoPathBin(cfg.GoPath))
		}
	}
	return globalBin, nil

}

var notFoundRe = regexp.MustCompile("evironment .* not found")

func isNotFound(err error) bool {
	err = errors.Cause(err)
	return notFoundRe.MatchString(err.Error())
}

func configGet(name string) (environment.Config, error) {
	basePath, err := paths.XdgDataConfig()
	if err != nil {
		return environment.Config{}, errors.Wrapf(err, "retrieving config for %q", name)
	}
	cfgs, err := environment.LoadConfig(basePath)
	if err != nil {
		return environment.Config{}, errors.Wrapf(err, "loading config for %q", name)
	}
	env, ok := cfgs[name]
	if !ok {
		return environment.Config{}, errors.Errorf("environment %q not found", name)
	}
	return env, nil
}

// Switch changes the environment to the specified one
// if it exists, otherwise its a noop and returns error.
func Switch(installName string) error {
	basePath, err := paths.XdgData()
	if err != nil {
		return errors.Wrapf(err, "retrieving config for %q", installName)
	}

	env, err := configGet(installName)
	if err != nil {
		return errors.Wrapf(err, "loading config to switch to %q", installName)
	}

	settings, err := environment.LoadSettings(basePath)
	if err != nil {
		return errors.Wrapf(err, "loading settings to switch to %q", installName)
	}

	extraBins, err := globalBins()
	if err != nil {
		return errors.Wrap(err, "determining global bin paths")
	}
	return errors.Wrapf(goswitch.Switch(env, installName == settings.Default, extraBins),
		"switching to environment %q", installName)
}

// Reset will try toreturn the env to its original state.
func Reset() error {
	return goswitch.Reset()
}

func ensureVersionInstalled(goVersion, goroot string) (goinstalls.Version, error) {
	goFolder, err := paths.XdgDataGoInstallsBinForVerson(goVersion)
	if err != nil {
		return goinstalls.Version{}, errors.Wrapf(err, "determining if %q exists", goVersion)
	}
	// if the bin folder is not there, the version is not properly
	// installed.
	if _, err := os.Stat(goFolder); err == nil {
		return goinstalls.VersionFromString(goVersion)
	}
	versions, err := goinstalls.OnlineAvailableVersions()
	if err != nil {
		return goinstalls.Version{}, errors.Wrap(err, "retrieving available online versions")
	}
	installFolder, err := paths.XdgDataGoInstalls()
	if err != nil {
		return goinstalls.Version{}, errors.Wrapf(err, "determining go installs folder to install %q", goVersion)
	}
	for k, v := range versions {
		if k.CommonVersionString() == goVersion {
			goFolder = filepath.Join(installFolder, k.String())
			if _, err := os.Stat(goFolder); err == nil {
				return k, nil
			}
			err = goinstalls.InstallVersion(k, v, installFolder, goroot)
			return k, errors.Wrapf(err, "installing go %q", goVersion)
		}
	}
	return goinstalls.Version{}, errors.Errorf("go version %q not found", goVersion)

}

// Create creates the an environment with the passed name
// in the passed go version, if it exists its a noop and
// returns an error.
func Create(installName, goVersion, goPath string, settings environment.Settings) error {
	_, err := configGet(installName)
	if err != nil && !isNotFound(err) {
		return errors.Wrapf(err, "determining if environment %q exists", installName)
	}
	v, err := ensureVersionInstalled(goVersion, settings.Goroot)
	if err != nil {
		return errors.Wrapf(err, "installing go %q to create %q environment", goVersion, installName)
	}
	c := environment.Config{
		Name:      installName,
		GoVersion: v.String(),
		GoPath:    goPath,
	}
	configPath, err := paths.XdgDataConfig()
	if err != nil {
		return errors.Wrapf(err, "getting config folder to save %q config", installName)
	}
	if err := c.Save(configPath); err != nil {
		return errors.Wrapf(err, "saving %q config", installName)
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
func Set(attribute, value string) error {
	environmentName, attribute, err := extractEnvironment(attribute)
	if err != nil {
		return errors.Errorf("setting %q to %q", attribute, value)
	}
	if environmentName != "" {
		cfg, err := configGet(environmentName)
		if err != nil {
			return errors.Wrapf(err, "finding config for %q", environmentName)
		}
		return errors.Wrapf(cfg.Set(attribute, value), "setting %q to %q", attribute, value)

		return errors.New("setting environment attributes is not implemented")
	}

	settingsFolder, err := paths.XdgData()
	if err != nil {
		return errors.Wrapf(err, "determining settings folder to set %q", attribute)
	}
	settings, err := environment.LoadSettings(settingsFolder)
	if err != nil {
		return errors.Wrap(err, "loading settings")
	}

	return errors.Wrapf(settings.Set(attribute, value), "setting %q to %q", attribute, value)
}
