package actions

import (
	"github.com/perrito666/goworkon/environment"
	"github.com/perrito666/goworkon/goswitch"
	"github.com/perrito666/goworkon/paths"
	"github.com/pkg/errors"
)

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
