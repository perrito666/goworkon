package environment

// Switch changes the environment to the specified one
// if it exists, otherwise its a noop and returns error.
func Switch(installName string) error {
	return nil
}

// Create creates the an environment with the passed name
// in the passed go version, if it exists its a noop and
// returns an error.
func Create(installName, goVersion string) error {
	return nil
}
