package environment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/juju/loggo"
	"github.com/pkg/errors"
)

var logger = loggo.GetLogger("goworkon.environment")

// SETTINGSFILE is the name of the file where the settings are stored.
const SETTINGSFILE = "settings.json"

// Config holds the information about a given environment.
type Config struct {
	// Name holds the name of the environment.
	Name string `json:"name"`
	// CompileSteps hold the commands to be run to compile this env main project.
	CompileSteps []string `json:"compilesteps"`
	// GoVersion holds the version of go this env should use.
	GoVersion string `json:"goversion"`
	// GlobalBin indicates if the $GOPATH/bin of this env will be added to PATH.
	GlobalBin bool `json:"globalbin"`
	// GoPath
	GoPath string `json:"gopath"`
}

func maybeEnsureFolderExists(folder string) error {
	fInfo, err := os.Stat(folder)
	if os.IsNotExist(err) {
		err = os.MkdirAll(folder, 0700)
		if err != nil {
			return errors.Wrapf(err, "creating %q", folder)
		}
		fInfo, err = os.Stat(folder)
	}
	if err != nil {
		return errors.Wrapf(err, "obtaining stat on %q", folder)
	}
	if !fInfo.IsDir() {
		return errors.Errorf("%q exists and it's not a folder", folder)
	}
	return nil
}

// Save serializes and writes the Config in a file in the
// passed folder.
func (c *Config) Save(baseFolder string) error {
	baseFolder = filepath.Join(baseFolder, "configs")
	if err := maybeEnsureFolderExists(baseFolder); err != nil {
		return errors.WithStack(err)
	}
	fileName := filepath.Join(baseFolder, fmt.Sprintf("%s.json", c.Name))
	fp, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
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
func LoadConfig(baseFolder string) (map[string]Config, error) {
	baseFolder = filepath.Join(baseFolder, "configs")
	if err := maybeEnsureFolderExists(baseFolder); err != nil {
		return nil, errors.WithStack(err)
	}
	files, err := filepath.Glob(filepath.Join(baseFolder, "*.json"))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	allConfigs := make(map[string]Config, len(files))
	for _, fileName := range files {
		err = func() error {
			var c Config
			contents, err := ioutil.ReadFile(fileName)
			if err != nil {
				return errors.WithStack(err)
			}
			if err = json.Unmarshal(contents, &c); err != nil {
				return errors.WithStack(err)
			}
			allConfigs[c.Name] = c
			return nil
		}()
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return allConfigs, nil
}

// Settings holds general settings for goworkon.
type Settings struct {
	// Goroot is the path to a working goroot, it is
	// required by the go compiler.
	// TODO (perrito666) make this updateable by each new version.
	Goroot string `json:"goroot"`
	// Default is the default environment to set, this will behave
	// a bit differently since its for general use.
	Default string `json:"default"`
}

// Save serializes and writes the Settings in a file in the
// passed folder.
func (s Settings) Save(baseFolder string) error {
	if err := maybeEnsureFolderExists(baseFolder); err != nil {
		return errors.WithStack(err)
	}
	fileName := filepath.Join(baseFolder, SETTINGSFILE)
	fp, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.Wrapf(err, "opening settings file %q for writing", fileName)
	}
	defer fp.Close()
	marshaled, err := json.Marshal(s)
	if err != nil {
		return errors.Wrap(err, "marshaling settings")
	}
	_, err = fp.Write(marshaled)
	if err != nil {
		return errors.Wrap(err, "writing marshaled settings")
	}
	return nil
}

// LoadSettings will load Settings files in the given location
func LoadSettings(baseFolder string) (Settings, error) {
	settingsFile := filepath.Join(baseFolder, SETTINGSFILE)
	logger.Debugf("loading settings from %q", settingsFile)
	_, err := os.Stat(settingsFile)
	if os.IsNotExist(err) {
		return Settings{}, nil
	}
	if err != nil {
		return Settings{}, errors.Wrap(err, "reading settings")
	}

	var s Settings
	contents, err := ioutil.ReadFile(settingsFile)
	if err != nil {
		return Settings{}, errors.WithStack(err)
	}
	if err := json.Unmarshal(contents, &s); err != nil {
		return Settings{}, errors.WithStack(err)
	}
	return s, nil
}
