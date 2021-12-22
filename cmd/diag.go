package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
)

var diagCmd = &cli.Command{
	Name:     "diag",
	Usage:    "diag the collected infos",
	Category: cmdDiag,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"C"},
			Usage:   "--config or -C",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"N"},
			Usage:   "node name",
			Value:   "",
		},
		&cli.StringSliceFlag{
			Name:    "diags",
			Aliases: []string{"D"},
			Usage:   "diags to analyze",
			Value:   cli.NewStringSlice(string(config.NoDiag)),
		},
	},
	Action: func(ctx *cli.Context) error {
		configPath := ctx.String("config")
		if configPath != "" {
			GlobalConfig, err := config.NewConfig(configPath, config.GetConfigType(configPath))
			if err != nil {
				return err
			}
			fmt.Printf("%+v", GlobalConfig.String())
		}
		if GlobalConfig == nil {
			return ErrConfigIsNull
		}
		if ctx.IsSet("name") {
			name := ctx.String("name")
			node := GlobalConfig.Nodes[name]
			if ctx.IsSet("diags") {
				diags := ctx.StringSlice("diags")
				diagOptions := make([]config.DiagOption, len(diags))
				for i := range diags {
					diagOptions[i] = config.DiagOption(diags[i])
				}
				node.Diags = diagOptions
			}
			fmt.Printf("%+v", GlobalConfig.String())
		}
		return nil
	},
}
