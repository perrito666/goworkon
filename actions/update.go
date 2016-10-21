package actions

import (
	"fmt"

	"github.com/perrito666/goworkon/environment"
	"github.com/perrito666/goworkon/goinstalls"
	"github.com/perrito666/goworkon/paths"
	"github.com/pkg/errors"
)

func update(v goinstalls.Version, envs []string) error {
	cfgData, err := paths.XdgDataConfig()
	if err != nil {
		return errors.Wrap(err, "getting path for config files")
	}
	for _, env := range envs {
		cfg, err := configGet(env)
		if err != nil {
			return errors.Wrapf(err, "updating %q to %q", env, v.String)
		}
		cfg.GoVersion = v.String()
		cfg.Save(cfgData)
		for _, step := range cfg.CompileSteps {
			fmt.Println(step)
			// TODO(perrito666) switch and run compile steps
		}
	}
	return nil
}

// UpdateToVersion updates the given environment to the given version or
// returns an error if not possible
func UpdateToVersion(environmentName string, version goinstalls.Version) error {
	fmt.Printf("will update %q to %q\n", environmentName, version.String())
	_, err := configGet(environmentName)
	if err != nil {
		return errors.Wrapf(err, "determining if environment %q exists", environmentName)
	}
	if err := ensureCanUpdateTo(version); err != nil {
		return errors.Wrapf(err, "installing go %q to update %q environment", version.String(), environmentName)
	}

	return errors.Wrapf(update(version, []string{environmentName}), "updating %q to version %q", environmentName, version.String())
}

func matchingVersion(v goinstalls.Version, haystack map[goinstalls.Version]string) (goinstalls.Version, bool) {
	for k := range haystack {
		if v.CommonVersionString() == k.CommonVersionString() {
			return k, true
		}
	}
	return goinstalls.Version{}, false
}

// UpdateAllTo will update all environments that share the common version
// to the passed patch.
func UpdateAllTo(version goinstalls.Version) error {
	basePath, err := paths.XdgDataConfig()
	if err != nil {
		return errors.Wrap(err, "retrieving configs for listing")
	}
	cfgs, err := environment.LoadConfig(basePath)
	if err != nil {
		return errors.Wrap(err, "loading configis for listing")
	}
	updateables := []string{}
	for _, cfg := range cfgs {
		v, err := goinstalls.VersionFromString(cfg.GoVersion)
		if err != nil {
			return errors.Wrapf(err, "finding version for %q", cfg.Name)
		}
		if version.CommonVersionString() == v.CommonVersionString() {
			updateables = append(updateables, cfg.Name)
		}
	}
	if len(updateables) > 0 {
		if version.Patch == 0 {
			versions, err := goinstalls.OnlineAvailableVersions()
			if err != nil {
				return errors.Wrap(err, "fetching available versions for udpate")
			}
			var ok bool
			version, ok = matchingVersion(version, versions)
			if !ok {
				return errors.Errorf("unavailable version %q", version.String())
			}
		}
		if err := ensureCanUpdateTo(version); err != nil {
			return errors.Wrapf(err, "installing go %q to update environments", version)
		}

		return errors.Wrapf(update(version, updateables), "updating all versions %q", version.String())
	}
	return nil

}

// UpdateToLatest will update the environment to the latest available version
// of go
func UpdateToLatest(environmentName string) error {
	version, _, err := goinstalls.NewestAvailableOnline()
	if err != nil {
		return errors.Wrap(err, "obtaining latest go version")
	}
	if err := ensureCanUpdateTo(version); err != nil {
		return errors.Wrapf(err, "installing go %q to update environment %q", version.String, environmentName)
	}

	return errors.Wrapf(update(version, []string{environmentName}), "updating %q to version %q", environmentName, version.String())
}
