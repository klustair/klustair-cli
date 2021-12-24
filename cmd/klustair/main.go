package main

import (
	_ "embed"
	"log"
	"os"

	"github.com/klustair/klustair-go/cmd/klustair/app"
)

//go:embed VERSION
var version string

func main() {
	cli := app.NewApp(version)
	err := cli.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
