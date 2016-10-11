package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/perrito666/goworkon/environment"
	"github.com/perrito666/goworkon/goinstalls"
)

// TODO(perrito666) this should be auto discovered from go repo.
const currentGoVersion = "1.7"

const (
	// COMMANDSWITCH is the name of the switch-to-env command.
	COMMANDSWITCH = "switch"
	// COMMANDCREATE is the name of the create-env command.
	COMMANDCREATE = "create"
	// COMMANDUPDATE is the name of the update-go-version command.
	COMMANDUPDATE = "update"
)

var (
	// args
	command  string
	switchTo string
	create   string

	// flags
	goVersion   string
	updateRepos bool
)

func init() {

	flag.StringVar(&goVersion, "go-version", currentGoVersion, "the go version to be used (if none specified, all be updated)")
	flag.BoolVar(&updateRepos, "update-envs", false, "update all envs that use this major go version")

	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("no command specified")
		os.Exit(1)
	}
	args := flag.Args()
	command = flag.Arg(0)
	switch flag.Arg(0) {
	case COMMANDSWITCH:
		if flag.NArg() != 2 {
			fmt.Println("unexpected number of arguments %d", flag.NArg())
			fmt.Println("the expected format is: goworkon switch <envname>")
			os.Exit(1)
		}
		switchTo = flag.Arg(1)
	case COMMANDCREATE:
		if flag.NArg() != 2 {
			fmt.Println("unexpected number of arguments %d", flag.NArg())
			fmt.Println("the expected format is: goworkon create <envname>")
			os.Exit(1)
		}

		create = flag.Arg(1)
	case COMMANDUPDATE:
		if flag.NArg() != 0 {
			fmt.Println("unexpected number of arguments %d", flag.NArg())
			fmt.Println("the expected format is: goworkon update")
			os.Exit(1)
		}
	default:
		fmt.Println("command %q is not supported", args[0])
		os.Exit(1)
	}

}

func main() {
	a, b := goinstalls.OnlineAvailableVersions()
	fmt.Println(a)
	fmt.Println(b)
	os.Exit(0)
	dataDir, err := xdgDataConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	configs, err := environment.LoadConfig(dataDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(configs)
	switch command {
	case COMMANDSWITCH:
		environment.Switch(switchTo)
	case COMMANDCREATE:
		environment.Create(create, goVersion)
	case COMMANDUPDATE:
		goinstalls.Update(goVersion, updateRepos)
	}
}
