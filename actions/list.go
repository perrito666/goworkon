package actions

import (
	"fmt"

	"github.com/perrito666/goworkon/environment"
	"github.com/perrito666/goworkon/paths"
	"github.com/pkg/errors"
)

// List prints a list of the existing configs.
func List() error {
	basePath, err := paths.XdgDataConfig()
	if err != nil {
		return errors.Wrap(err, "retrieving configs for listing")
	}
	cfgs, err := environment.LoadConfig(basePath)
	if err != nil {
		return errors.Wrap(err, "loading configis for listing")
	}
	for _, cfg := range cfgs {
		fmt.Println(fmt.Sprintf("(%s) %q:%s", cfg.GoVersion, cfg.Name, cfg.GoPath))
		if len(cfg.CompileSteps) > 0 {
			for i, step := range cfg.CompileSteps {
				fmt.Println("_%d: %q", i, step)
			}
		}
	}
	return nil
}
