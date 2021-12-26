package cmd

import (
	"github.com/urfave/cli/v2"
	"github.com/xjlgod/nebula-database-diagnostic/intrenal/info"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"strings"
)

var infoCmd = &cli.Command{
	Name:     "info",
	Usage:    "fetch the nebula graph infos",
	Category: cmdDiag,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"C"},
			Usage:   "--config or -C, the config file for fetching infos",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"N"},
			Usage:   "--node or -N, to modify the info config by name",
			Value:   "",
		},
		&cli.StringFlag{
			Name:    "output_dir_path",
			Aliases: []string{"O"},
			Usage:   "--output_dir_path or -O, the output dir of infos, logs, and others output",
			Value:   "./output",
		},
		&cli.BoolFlag{
			Name:    "log_to_file",
			Aliases: []string{"L"},
			Usage:   "--log_to_file or -L, log to file or to cmd",
			Value:   false,
		},
		&cli.StringFlag{
			Name:    "duration",
			Aliases: []string{"D"},
			Usage:   "--duration or -D, the fetch duration, '-1' means that fetch the infos with infinity",
			Value:   "-1",
		},
		&cli.StringFlag{
			Name:    "period",
			Aliases: []string{"P"},
			Usage:   "--period or -P, the period when fetch the infos",
			Value:   "5s",
		},
		&cli.StringFlag{
			Name:    "infos",
			Aliases: []string{"I"},
			Usage:   "--infos or -I, the infos to fetch, will overwrite the infos in config, included: metrics, service, physical, all, no, etc.",
			Value:   string(config.AllInfo),
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
			inf := GlobalConfig.Infos[name]
			if ctx.IsSet("output_dir_path") {
				outputDirPath := ctx.String("output_dir_path")
				inf.Output.DirPath = outputDirPath
			}
			if ctx.IsSet("log_to_file") {
				logToFile := ctx.Bool("log_to_file")
				inf.Output.LogToFile = logToFile
			}
			if ctx.IsSet("duration") {
				duration := ctx.String("duration")
				inf.Duration = duration
			}
			if ctx.IsSet("period") {
				period := ctx.String("period")
				inf.Period = period
			}
			if ctx.IsSet("infos") {
				var infos []string
				infosStr := ctx.String("infos")
				if strings.Contains(infosStr, "all") {
					infos = []string{"all"}
				} else {
					for _, infoStr := range strings.Split(infosStr, ",") {
						infos = append(infos, infoStr)
					}
				}

				infoOptions := make([]config.InfoOption, len(infos))
				for i := range infos {
					infoOptions[i] = config.InfoOption(infos[i])
				}
				inf.Options = infoOptions
			}
		}
		info.Run(GlobalConfig)
		return nil
	},
}
