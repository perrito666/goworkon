package goswitch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/perrito666/goworkon/environment"
	"github.com/perrito666/goworkon/paths"
	"github.com/pkg/errors"
)

const (
	// PATH is the name of the env variable of the same name.
	PATH = "PATH"
	// GOPATH is the name of the env variable of the same name.
	GOPATH = "GOPATH"
	// PS1 is the name of the env variable of the same name.
	PS1 = "PS1"
	// CDPATH is the name of the env variable of the same name.
	CDPATH = "CDPATH"
	// PREVPATH is the name of the variable used to backup PATH.
	PREVPATH = "GOWORKON_PREVIOUS_PATH"
	// PREVGOPATH is the name of the variable used to backup GOPATH.
	PREVGOPATH = "GOWORKON_PREVIOUS_GOPATH"
	// PREVPS1 is the name of the variable used to backup PS1.
	PREVPS1 = "GOWORKON_PREVIOUS_PS1"
	// PREVCDPATH is the name of the variable used to backup CDPATH.
	PREVCDPATH = "GOWORKON_PREVIOUS_CDPATH"
)

func setenv(varName, varValue string) string {
	return fmt.Sprintf("%s=%s", varName, varValue)
}

// Switch will set the proper environment variables to set an environment
// as the current running one.
func Switch(cfg environment.Config, isDefault bool, extraBin []string) error {
	pgopath := os.Getenv(PREVGOPATH)
	ppath := os.Getenv(PREVPATH)
	pps1 := os.Getenv(PREVPS1)
	pcdpath := os.Getenv(PREVCDPATH)
	gopath := os.Getenv(GOPATH)
	path := os.Getenv(PATH)
	ps1 := os.Getenv(PS1)
	cdpath := os.Getenv(CDPATH)

	envVars := []string{}
	// backup vanilla paths.
	if pgopath == "" && !isDefault {
		envVars = append(envVars, setenv(PREVGOPATH, gopath))
	}
	if ppath == "" && !isDefault {
		envVars = append(envVars, setenv(PREVPATH, path))
	}
	if pps1 == "" && !isDefault {
		envVars = append(envVars, setenv(PREVPS1, fmt.Sprintf("\"%s\"", ps1)))
	}
	if pcdpath == "" {
		envVars = append(envVars, setenv(PREVCDPATH, cdpath))
	}
	// set env vars.
	envVars = append(envVars, setenv(GOPATH, cfg.GoPath))
	newCdpath := filepath.Join(cfg.GoPath, "src")
	if cdpath != "" && pcdpath == "" {
		newCdpath = fmt.Sprintf("%s:%s", cdpath, newCdpath)
	}
	envVars = append(envVars, setenv(CDPATH, newCdpath))
	if len(extraBin) > 0 {
		extraBin = append(extraBin, path)
		path = strings.Join(extraBin, ":")
	}
	goInstallsPath, err := paths.XdgDataGoInstallsBinForVerson(cfg.GoVersion)
	if err != nil {
		return errors.Wrapf(err, "trying to determine go installs path to switch to %q", cfg.Name)
	}
	newPath := paths.PATHInsert(path, paths.GoPathBin(cfg.GoPath), goInstallsPath)
	envVars = append(envVars, setenv(PATH, newPath))
	// Default env does not need new ps1
	if pps1 != "" {
		ps1 = pps1
	}
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
