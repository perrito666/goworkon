package main

import (
	"github.com/perrito666/goworkon/actions"
	"github.com/perrito666/goworkon/environment"
	"github.com/perrito666/goworkon/goinstalls"
	"github.com/pkg/errors"
)

// Command interface represents a command of goworkon.
type Command interface {
	// Usage returns a string explaining the usage of this command.
	Usage() string
	// Validate returns error if some of the pre-requisites for
	// the command are not met.
	Validate() error
	// Run executes the command and returns error if there was a problem
	Run() error
}

// Create command allows for the creation of environments.
type Create struct {
	environmentName string
	goVersion       string
	goPath          string
	settings        environment.Settings
}

// Usage implements Command.
func (c Create) Usage() string {
	return "the expected format is: goworkon [Options] create <envname>"
}

// Validate implements Command.
func (c Create) Validate() error {
	if c.environmentName == "" {
		return errors.New("missing environment name")
	}
	if c.goVersion == "" {
		return errors.New("missing go version")
	}
	if c.goPath == "" {
		return errors.New("missing gopath/workspace for the environment")
	}
	return nil
}

// Run implements Command.
func (c Create) Run() error {
	err := actions.Create(c.environmentName, c.goVersion, c.goPath, c.settings)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil

}

// Switch command allows the user to change the working environment
// to a specified one.
type Switch struct {
	environmentName string
}

// Usage implements Command.
func (s Switch) Usage() string {
	return "the expected format is: goworkon [Options] switch [envname]\n" +
		"if <envname> is not provided, switch will reset to default"
}

// Validate implements Command.
func (s Switch) Validate() error {
	return nil
}

// Run implements Command.
func (s Switch) Run() error {
	if s.environmentName == "" {
		return actions.Reset()
	}
	err := actions.Switch(s.environmentName)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

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
	if u.environmentName == "" && goVersion == "" {
		return errors.New("specify either a go version or an environment name")
	}
	v, err := goinstalls.VersionFromString(goVersion)
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
	return nil
}

// Set command helps the user set configuration values for environments and
// global settings.
type Set struct {
	attribute string
	value     string
}

// Usage implements Command.
func (s Set) Usage() string {
	// TODO(perrito666) remember to forbid spaces on env names.
	return "the expected format is: goworkon set attribute[@environment] value"
}

// Validate implements Command.
func (s Set) Validate() error {
	if s.attribute == "" {
		return errors.New("the attribute cannot be empty")
	}
	return nil
}

// Run implements Command.
func (s Set) Run() error {
	return errors.WithStack(actions.Set(s.attribute, s.value))
}
