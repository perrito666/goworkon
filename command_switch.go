package main

import (
	"github.com/perrito666/goworkon/actions"
	"github.com/pkg/errors"
)

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
