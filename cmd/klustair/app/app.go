package app

import (
	"github.com/urfave/cli/v2"
)

var (
	namespacesFlag = cli.StringFlag{
		Name:    "namespaces",
		Aliases: []string{"n"},
		Value:   "",
		Usage:   "Coma separated whitelist of Namespaces to check",
		EnvVars: []string{"KLUSTAIR_NAMESPACES"},
	}
	namespacesblacklistFlag = cli.StringFlag{
		Name:    "namespacesblacklist",
		Aliases: []string{"N"},
		Value:   "",
		Usage:   "Coma separated whitelist of Namespaces to check",
		EnvVars: []string{"KLUSTAIR_NAMESPACESBLACKLIST"},
	}

	kubeauditFlag = cli.StringFlag{
		Name:    "kubeaudit",
		Aliases: []string{"k"},
		Value:   "all",
		Usage:   "Coma separated list of audits to run. (disable: \"none\")",
		EnvVars: []string{"KLUSTAIR_KUBEAUDIT"},
	}

	trivyFlag = cli.BoolFlag{
		Name:    "trivy",
		Aliases: []string{"t"},
		Usage:   "Run Trivy vulnerability checks",
		EnvVars: []string{"KLUSTAIR_TRIVY"},
	}

	labelFlag = cli.StringFlag{
		Name:    "label",
		Aliases: []string{"l"},
		Value:   "",
		Usage:   "A optional title for your run",
		EnvVars: []string{"KLUSTAIR_NAMESPACESBLACKLIST"},
	}

	trivycredentialspathFlag = cli.StringFlag{
		Name:    "repocredentialspath",
		Aliases: []string{"c"},
		Value:   "",
		Usage:   "Path to repo credentials for trivy",
		EnvVars: []string{"KLUSTAIR_REPOCREDENTIALSPATH"},
	}

	limitdateFlag = cli.IntFlag{
		Name:    "limitdate",
		Aliases: []string{"ld"},
		Value:   0,
		Usage:   "Remove reports older than X days",
		EnvVars: []string{"KLUSTAIR_LIMITDATE"},
	}

	limitnrFlag = cli.IntFlag{
		Name:    "limitnr",
		Aliases: []string{"ln"},
		Value:   0,
		Usage:   "Keep only X reports",
		EnvVars: []string{"KLUSTAIR_LIMITNR"},
	}

	configkeyFlag = cli.StringFlag{
		Name:    "configkey",
		Aliases: []string{"C"},
		Value:   "",
		Usage:   "Load remote configuration from frontend",
		EnvVars: []string{"KLUSTAIR_CONFIGKEY"},
	}

	debugFlag = cli.BoolFlag{
		Name:    "debug",
		Aliases: []string{"d"},
		Usage:   "debug mode",
		EnvVars: []string{"KLUSTAIR_DEBUG"},
	}

	quietFlag = cli.BoolFlag{
		Name:    "verbose",
		Aliases: []string{"V"},
		Usage:   "increase output verbosity",
		EnvVars: []string{"KLUSTAIR_VERBOSE"},
	}

	// Global flags
	globalFlags = []cli.Flag{
		&quietFlag,
		&debugFlag,
	}

	imageFlags = []cli.Flag{
		&namespacesFlag,
		&namespacesblacklistFlag,
		&kubeauditFlag,
		&trivyFlag,
		&labelFlag,
		&trivycredentialspathFlag,
		&limitdateFlag,
		&limitnrFlag,
		&configkeyFlag,
	}
)

// NewApp is the factory method to return Trivy CLI
func NewApp(version string) *cli.App {

	app := cli.NewApp()
	app.Name = "klustair"
	app.Version = version
	app.Usage = "A simple and comprehensive vulnerability scanner for kubernetes"
	app.EnableBashCompletion = true

	flags := append(globalFlags, imageFlags...)
	app.Flags = flags

	return app
}
