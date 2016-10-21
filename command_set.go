package main

import (
	"github.com/perrito666/goworkon/actions"
	"github.com/pkg/errors"
)

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
