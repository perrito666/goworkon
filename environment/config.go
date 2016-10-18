package environment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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

// Set will set the value of <attribute> to <value> if attribute is a valid
// member of Config.
func (c Config) Set(attribute, value, baseFolder string) error {
	switch strings.ToLower(attribute) {
	case "globalbin":
		if strings.ToLower(value) == "true" {
			c.GlobalBin = true
		} else {
			c.GlobalBin = false
		}
	case "compilesteps":
		c.CompileSteps = strings.Split(value, ";")
	default:
		return errors.Errorf("%q is not a valid setting", attribute)
	}
	return errors.WithStack(c.Save(baseFolder))
}
