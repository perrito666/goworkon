package actions

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/perrito666/goworkon/environment"
	"github.com/perrito666/goworkon/goinstalls"
	"github.com/perrito666/goworkon/paths"
	"github.com/pkg/errors"
)

func globalBins() ([]string, error) {
	basePath, err := paths.XdgDataConfig()
	if err != nil {
		return nil, errors.Wrap(err, "retrieving config files")
	}
	cfgs, err := environment.LoadConfig(basePath)
	if err != nil {
		return nil, errors.Wrap(err, "loading configs")
	}
	globalBin := []string{}
	for _, cfg := range cfgs {
		if cfg.GlobalBin {
			globalBin = append(globalBin, paths.GoPathBin(cfg.GoPath))
		}
	}
	return globalBin, nil

}

var notFoundRe = regexp.MustCompile("environment .* not found")

func isNotFound(err error) bool {
	err = errors.Cause(err)
	return notFoundRe.MatchString(err.Error())
}

func configGet(name string) (environment.Config, error) {
	basePath, err := paths.XdgDataConfig()
	if err != nil {
		return environment.Config{}, errors.Wrapf(err, "retrieving config for %q", name)
	}
	cfgs, err := environment.LoadConfig(basePath)
	if err != nil {
		return environment.Config{}, errors.Wrapf(err, "loading config for %q", name)
	}
	env, ok := cfgs[name]
	if !ok {
		return environment.Config{}, errors.Errorf("environment %q not found", name)
	}
	return env, nil
}

func ensureCanUpdateTo(goVersion goinstalls.Version) error {
	dataDir, err := paths.XdgData()
	if err != nil {
		return errors.Wrap(err, "finding data dir")
	}

	settings, err := environment.LoadSettings(dataDir)
	if err != nil {
		return errors.Wrap(err, "loading settings")
	}

	_, err = ensureVersionInstalled(goVersion.String(), settings.Goroot)
	if err != nil {
		return errors.Wrapf(err, "installing go %q", goVersion)
	}
	return nil
}

func ensureVersionInstalled(goVersion, goroot string) (goinstalls.Version, error) {
	goFolder, err := paths.XdgDataGoInstallsBinForVerson(goVersion)
	if err != nil {
		return goinstalls.Version{}, errors.Wrapf(err, "determining if %q exists", goVersion)
	}
	// if the bin folder is not there, the version is not properly
	// installed.
	if _, err := os.Stat(goFolder); err == nil {
		return goinstalls.VersionFromString(goVersion)
	}
	versions, err := goinstalls.OnlineAvailableVersions()
	if err != nil {
		return goinstalls.Version{}, errors.Wrap(err, "retrieving available online versions")
	}
	installFolder, err := paths.XdgDataGoInstalls()
	if err != nil {
		return goinstalls.Version{}, errors.Wrapf(err, "determining go installs folder to install %q", goVersion)
	}
	reqVersion, err := goinstalls.VersionFromString(goVersion)
	if err != nil {
		return goinstalls.Version{}, errors.Wrapf(err, "parsing the requested version %q", goVersion)
	}
	for k, v := range versions {
		minorMatch := reqVersion.Patch == 0 && k.CommonVersionString() == goVersion
		fullMatch := reqVersion.SameVersion(k)
		if minorMatch || fullMatch {
			goFolder = filepath.Join(installFolder, k.String())
			if _, err := os.Stat(goFolder); err == nil {
				return k, nil
			}
			err = goinstalls.InstallVersion(k, v, installFolder, goroot)
			return k, errors.Wrapf(err, "installing go %q", goVersion)
		}
	}
	return goinstalls.Version{}, errors.Errorf("go version %q not found", goVersion)

}

func extractEnvironment(attribute string) (string, string, error) {
	parts := strings.Split(attribute, "@")
	l := len(parts)
	if l == 2 {
		return parts[0], parts[1], nil
	}
	if l == 1 {
		return "", parts[0], nil
	}
	return "", "", errors.Errorf("%q is not a valid attribute", attribute)
}
