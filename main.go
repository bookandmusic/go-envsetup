package main

import (
	"log"
	"os"

	"github.com/bookandmusic/envsetup/cli"
)

func main() {
	app := cli.CreateApp()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
