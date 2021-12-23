package cmd

import (
	"github.com/urfave/cli/v2"
	"github.com/xjlgod/nebula-database-diagnostic/intrenal/info"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
)

var infoCmd = &cli.Command{
	Name:     "info",
	Usage:    "fetch the graph infos",
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
		&cli.StringFlag{
			Name:    "output_dir_path",
			Aliases: []string{"O"},
			Usage:   "output dir",
			Value:   "./out",
		},
		&cli.BoolFlag{
			Name:    "output_remote",
			Aliases: []string{"R"},
			Usage:   "output dir on remote",
			Value:   false,
		},
		&cli.StringFlag{
			Name:    "duration",
			Aliases: []string{"D"},
			Usage:   "diag duration",
			Value:   "-1",
		},
		&cli.StringFlag{
			Name:    "period",
			Aliases: []string{"P"},
			Usage:   "info period",
			Value:   "5s",
		},
		&cli.StringSliceFlag{
			Name:    "infos",
			Aliases: []string{"I"},
			Usage:   "infos to fetch, will overwrite the infos in config",
			Value:   cli.NewStringSlice(string(config.AllInfo)),
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
		}
		if GlobalConfig == nil {
			return ErrConfigIsNull
		}
		if ctx.IsSet("name") {
			name := ctx.String("name")
			node := GlobalConfig.Nodes[name]
			if ctx.IsSet("output_dir_path") {
				outputDirPath := ctx.String("output_dir_path")
				node.Output.DirPath = outputDirPath
			}
			if ctx.IsSet("output_remote") {
				outputRemote := ctx.Bool("output_remote")
				node.Output.Remote = outputRemote
			}
			if ctx.IsSet("duration") {
				duration := ctx.String("duration")
				node.Duration = duration
			}
			if ctx.IsSet("period") {
				period := ctx.String("period")
				node.Period = period
			}
			if ctx.IsSet("infos") {
				infos := ctx.StringSlice("infos")
				infoOptions := make([]config.InfoOption, len(infos))
				for i := range infos {
					infoOptions[i] = config.InfoOption(infos[i])
				}
				node.Infos = infoOptions
			}
		}
		info.Run(GlobalConfig)
		return nil
	},
}
