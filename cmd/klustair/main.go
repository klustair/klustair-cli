package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/klustair/klustair-go/pkg/commands"
)

//go:embed VERSION
var version string

func main() {
	app := commands.NewApp(version)
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
