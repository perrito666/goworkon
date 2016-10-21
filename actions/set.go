package actions

import (
	"github.com/perrito666/goworkon/environment"
	"github.com/perrito666/goworkon/paths"
	"github.com/pkg/errors"
)

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
