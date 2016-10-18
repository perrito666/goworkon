package goswitch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/perrito666/goworkon/environment"
)

const (
	// PATH is the name of the env variable of the same name.
	PATH = "PATH"
	// GOPATH is the name of the env variable of the same name.
	GOPATH = "GOPATH"
	// PS1 is the name of the env variable of the same name.
	PS1 = "PS1"
	// PREVPATH is the name of the variable used to backup PATH.
	PREVPATH = "GOWORKON_PREVIOUS_PATH"
	// PREVGOPATH is the name of the variable used to backup GOPATH.
	PREVGOPATH = "GOWORKON_PREVIOUS_GOPATH"
	// PREVPS1 is the name of the variable used to backup PS1.
	PREVPS1 = "GOWORKON_PREVIOUS_PS1"
)

func setenv(varName, varValue string) string {
	return fmt.Sprintf("%s=%s", varName, varValue)
}

// Switch will set the proper environment variables to set an environment
// as the current running one.
func Switch(basePath string, cfg environment.Config, isDefault bool, extraBin []string) error {
	pgopath := os.Getenv(PREVGOPATH)
	ppath := os.Getenv(PREVPATH)
	pps1 := os.Getenv(PREVPS1)
	gopath := os.Getenv(GOPATH)
	path := os.Getenv(PATH)
	ps1 := os.Getenv(PS1)
	envVars := []string{}
	if pgopath == "" {
		envVars = append(envVars, setenv(PREVGOPATH, gopath))
	}
	envVars = append(envVars, setenv(GOPATH, cfg.GoPath))
	if ppath == "" {
		envVars = append(envVars, setenv(PREVPATH, path))
	}
	if pps1 == "" {
		envVars = append(envVars, setenv(PREVPS1, fmt.Sprintf("\"%s\"", ps1)))
	}
	// TODO(perrito) go throught the environment and remove the path section
	// that holds the current go install.
	// Also threat default as a special cookie and make it leave no trace.
	envVars = append(envVars, setenv(GOPATH, cfg.GoPath))
	if len(extraBin) > 0 {
		extraBin = append(extraBin, path)
		path = strings.Join(extraBin, ":")
	}
	newPath := strings.Join([]string{filepath.Join(cfg.GoPath, "bin"),
		// TODO(perrito) this implies heavy out of band knowledge, lets
		// store these things somewhere instead or build a unique source
		// of truth for paths.
		filepath.Join(basePath, cfg.GoVersion, "go", "bin"),
		path}, ":")
	envVars = append(envVars, setenv(PATH, newPath))
	// Default env will
	if !isDefault {
		envVars = append(envVars, setenv(PS1, fmt.Sprintf("\"%s(%s)$ \"", ps1, cfg.Name)))
	}
	fmt.Println(strings.Join(envVars, "\n"))
	return nil
}

// Reset will set the environment to its previous state.
func Reset() error {
	envVars := []string{
		setenv(PREVPATH, ""),
		setenv(PREVGOPATH, ""),
		setenv(PREVPS1, ""),
	}
	pgopath := os.Getenv(PREVGOPATH)
	ppath := os.Getenv(PREVPATH)
	envVars = append(envVars, setenv(GOPATH, pgopath))

	if ppath != "" {
		envVars = append(envVars, setenv(PATH, ppath))
	}

	pps1 := os.Getenv(PREVPS1)
	if pps1 != "" {
		envVars = append(envVars, setenv(PS1, fmt.Sprintf("\"%s\"", pps1)))
	}
	fmt.Println(strings.Join(envVars, "\n"))
	return nil
}
