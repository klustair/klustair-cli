package commands

import "github.com/urfave/cli"

// NewApp is the factory method to return Trivy CLI
func NewApp(version string) *cli.App {

	app := cli.NewApp()
	app.Name = "klustair"
	app.Version = version
	app.ArgsUsage = "target"
	app.Usage = "A simple and comprehensive vulnerability scanner for containers"
	app.EnableBashCompletion = true
	return app
}
