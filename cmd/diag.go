package cmd

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/xjlgod/nebula-database-diagnostic/intrenal/diag"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"strings"
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
			Name:    "output_dir_path",
			Aliases: []string{"O"},
			Usage:   "output dir",
			Value:   "./output",
		},
		&cli.BoolFlag{
			Name:    "log_to_file",
			Aliases: []string{"L"},
			Usage:   "log to file or to cmd",
			Value:   false,
		},
		&cli.StringFlag{
			Name:    "input_dir_path",
			Aliases: []string{"I"},
			Usage:   "input dir",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "diags",
			Aliases: []string{"D"},
			Usage:   "diags to analyze",
			Value:   string(config.NoDiag),
		},
	},
	Action: func(ctx *cli.Context) error {
		configPath := ctx.String("config")
		var err error
		if configPath != "" {
			GlobalConfig, err = config.NewConfig(configPath, config.GetConfigType(configPath))
			if err != nil {
				return err
			}
			fmt.Printf("%+v", GlobalConfig.String())
		}
		if GlobalConfig == nil {
			GlobalConfig = new(config.Config)
			config.ConfigComplete(GlobalConfig)
		}

		dia := GlobalConfig.Diag
		if ctx.IsSet("output_dir_path") {
			outputDirPath := ctx.String("output_dir_path")
			dia.Output.DirPath = outputDirPath
		}
		if ctx.IsSet("log_to_file") {
			logToFile := ctx.Bool("log_to_file")
			dia.Output.LogToFile = logToFile
		}
		if ctx.IsSet("input_dir_path") {
			inputDirPath := ctx.String("input_dir_path")
			dia.Input.DirPath = inputDirPath
		} else {
			return ErrNoInput
		}
		if ctx.IsSet("diags") {
			var diags []string
			diagsStr := ctx.String("diags")
			if strings.Contains(diagsStr, "all") {
				diags = []string{"all"}
			} else {
				for _, diagStr := range strings.Split(diagsStr, ",") {
					diags = append(diags, diagStr)
				}
			}

			diagOptions := make([]config.DiagOption, len(diags))
			for i := range diags {
				diagOptions[i] = config.DiagOption(diags[i])
			}
			dia.Options = diagOptions
		}
		diag.Run(GlobalConfig)
		return nil
	},
}
