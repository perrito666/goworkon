package actions

import (
	"github.com/perrito666/goworkon/environment"
	"github.com/perrito666/goworkon/paths"
	"github.com/pkg/errors"
)

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
