package main

import (
	"log"
	"os"

	"github.com/klustair/klustair-go/pkg/commands"
)

var (
	version = "dev"
)

func main() {
	app := commands.NewApp(version)
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
