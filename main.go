package main

import (
	"github.com/urfave/cli/v2"
	"github.com/xjlgod/nebula-database-diagnostic/cmd"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:                   cmd.Name,
		Usage:                  cmd.Desc,
		Version:                cmd.Version,
		UseShortOptionHandling: true,
		Flags:                  cmd.GlobalOptions,
		Before:                 cmd.LoadGlobalOptions,
		Commands:               cmd.Commands,
	}
	err := app.Run(os.Args)
	if err != nil && err != cmd.ErrPrintAndExit {
		log.Fatal(err)
	}
}
