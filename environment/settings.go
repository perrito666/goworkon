package environment

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Settings holds general settings for goworkon.
type Settings struct {
	// Goroot is the path to a working goroot, it is
	// required by the go compiler.
	// TODO (perrito666) make this updateable by each new version.
	Goroot string `json:"goroot"`
	// Default is the default environment to set, this will behave
	// a bit differently since its for general use.
	Default string `json:"default"`

	// filePath holds the path for this settings file.
	filePath string
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

// Set will set the value of <attribute> to <value> if attribute is a valid
// member of Settings.
func (s Settings) Set(attribute, value string) error {
	if s.filePath == "" {
		return errors.New("these settings neds to be saved before Set can be used")
	}

	switch strings.ToLower(attribute) {
	case "goroot":
		s.Goroot = value
	case "default":
		s.Default = value
	default:
		return errors.Errorf("%q is not a valid setting", attribute)
	}
	return errors.WithStack(s.Save(s.filePath))
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
	s.filePath = baseFolder
	return s, nil
}
