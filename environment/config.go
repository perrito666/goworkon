package environment

import (
	"encoding/json"
	"os"
	"os/path"
	"path/filepath"

	"github.com/pkg/errors"
)

// Config holds the information about a given environment.
type Config struct {
	// Name holds the name of the environment.
	Name string `json:"Name"`
	// CompileSteps hold the commands to be run to compile this env main project.
	CompileSteps []string `json:"CompileSteps"`
	// GoVersion holds the version of go this env should use.
	GoVersion string `json:"GoVersion"`
	// GlobalBin indicates if the $GOPATH/bin of this env will be added to PATH.
	GlobalBin bool `json:"GlobalBin"`
}

func maybeEnsureFolderExists(folder string) error {
	fInfo, err := os.Stat(folder)
	if os.IsNotExists(err) {
		err = path.MkDirAll(folder, 0600)
		if err != nil {
			return errors.Wrapf(err, "creating %q", folder)
		}
		fInfo, err = os.Stat(folder)
	}
	if err != nil {
		return errors.Wrapf(err, "obtaining stat on %q", folder)
	}
	if !fInfo.IsDir() {
		return errors.New("%q exists and it's not a folder", folder)
	}
	return nil
}

// Save serializes and writes the Config in a file in the
// passed folder.
func (c *Config) Save(baseFolder string) error {
	if err := maybeEnsureFolderExists(baseFolder); err != nil {
		return errors.WithStack(err)
	}
	fileName := filepath.Join(baseFolder, c.Name)
	fp, err := os.Open(fileName)
	if err != nil {
		return errors.Wrapf(err, "opening config file %q for writing", fileName)
	}
	defer fp.Close()
	marshaled, err := json.Marshal(c)
	if err != nil {
		return errors.Wrapf(err, "marshaling config for %q", c.Name)
	}
	_, err = fp.Write(marshaled)
	if err != nil {
		return errors.Wrap(err, "writing marshaled data")
	}
	return nil
}

// LoadConfig will load Config files in the given location
func LoadConfig(baseFolder string) ([]Config, error) {
	return []Config{}, nil
}
