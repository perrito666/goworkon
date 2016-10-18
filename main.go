package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/juju/loggo"
	flag "github.com/ogier/pflag"
	"github.com/perrito666/goworkon/environment"
	"github.com/pkg/errors"
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
	// COMMANDSET is the name of the set attribute command.
	COMMANDSET = "set"
)

var (
	// flags
	goVersion string
)

var logger = loggo.GetLogger("goworkon")

func init() {
	//loggo.ConfigureLoggers(`<root>=DEBUG`)
	flag.StringVar(&goVersion, "go-version", currentGoVersion, "the go version to be used (if none specified, all be updated)")
}

func checkCommand(s environment.Settings) (Command, error) {
	if flag.NArg() == 0 {
		return nil, errors.New("no command specified")
	}
	if goVersion == "" {
		goVersion = currentGoVersion
	}

	switch flag.Arg(0) {
	case COMMANDSWITCH:
		return Switch{
			environmentName: flag.Arg(1),
		}, nil
	case COMMANDCREATE:
		return Create{
			environmentName: flag.Arg(1),
			goPath:          flag.Arg(2),
			goVersion:       goVersion,
			settings:        s,
		}, nil
	case COMMANDUPDATE:
		return Update{
			environmentName: flag.Arg(1),
			goVersion:       goVersion,
		}, nil
	case COMMANDSET:
		return Set{
			attribute: flag.Arg(1),
			value:     flag.Arg(2),
		}, nil
	}
	return nil, errors.Errorf("Unknown command %q\n", flag.Arg(0))
}

func promptData(query string) (string, error) {
	stdin := bufio.NewReader(os.Stdin)
	fmt.Print(query)
	answer, err := stdin.ReadString('\n')
	if err != nil {
		return "", errors.WithStack(err)
	}
	return strings.Trim(answer, "\n"), nil
}

func main() {
	var err error
	fail := func() {
		// TODO (perrito) make the + on format optional
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}

	flag.Parse()

	dataDir, err := xdgData()
	if err != nil {
		fail()
	}

	settings, err := environment.LoadSettings(dataDir)
	if err != nil {
		fail()
	}
	if settings.Goroot == "" {
		settings.Goroot, err = promptData("Please provide a valid GOROOT path: ")
		if err != nil {
			fail()
		}
		settings.Save(dataDir)
	}

	c, err := checkCommand(settings)
	if err != nil {
		fail()
	}

	if err = c.Validate(); err != nil {
		fmt.Printf("%v\n", err)
		fmt.Println(c.Usage())
		os.Exit(1)
	}

	err = c.Run()
	if err != nil {
		fail()
	}
}
