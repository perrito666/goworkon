package main

import (
	"github.com/perrito666/goworkon/actions"
	"github.com/perrito666/goworkon/goinstalls"
	"github.com/pkg/errors"
)

// Update command allows the user to change the go version of one or multiple
// environments to the latest.
type Update struct {
	environmentName string
	goVersion       string
}

// Usage implements Command.
func (u Update) Usage() string {
	return "the expected format is: goworkon [Options] switch <envname>"
}

// Validate implements Command.
func (u Update) Validate() error {
	if u.environmentName == "" && u.goVersion == "" {
		return errors.New("specify either a go version or an environment name")
	}
	v, err := goinstalls.VersionFromString(u.goVersion)
	if err != nil {
		return errors.WithStack(err)
	}
	// if only version is passed, version should be just a minor, this implies that
	// we will update all x.y installs to the latest x.y
	if u.environmentName == "" && v.Patch != 0 {
		return errors.New("when passing only go version, ommit the patch version (major.minor.patch)")
	}

	return nil
}

// Run implements Command.
func (u Update) Run() error {
	v, err := goinstalls.VersionFromString(u.goVersion)
	if err != nil {
		return errors.WithStack(err)
	}

	if u.environmentName == "" {
		return errors.WithStack(actions.UpdateAllTo(v))
	}
	if u.goVersion != "" {
		return errors.WithStack(actions.UpdateToVersion(u.environmentName, v))
	}
	return errors.WithStack(actions.UpdateToLatest(u.environmentName))
}
