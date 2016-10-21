package main

import (
	"github.com/perrito666/goworkon/actions"
	"github.com/pkg/errors"
)

// List commamd prints a list of all the existing environments.
type List struct {
}

// Usage implements Command.
func (l List) Usage() string {
	return "the expected format is: goworkon list"
}

// Validate implements Command.
func (l List) Validate() error {
	return nil
}

// Run implements Command.
func (l List) Run() error {
	return errors.WithStack(actions.List())
}
