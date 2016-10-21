package main

import (
	"github.com/perrito666/goworkon/actions"
	"github.com/perrito666/goworkon/environment"
	"github.com/pkg/errors"
)

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
