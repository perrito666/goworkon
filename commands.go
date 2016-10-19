package main

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
