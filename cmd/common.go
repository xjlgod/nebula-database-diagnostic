package cmd

import (
	"errors"
	"github.com/urfave/cli/v2"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
)

const (
	Name    = "nebula-diag-cli"
	Desc    = `A free and open source distributed diagnostic tool for nebula graph`
	Version = "v0.0.1"
)

const (
	cmdDiag = "Diag function commands"
)

var (
	ErrPrintAndExit = errors.New("print and exit")
	ErrConfigIsNull = errors.New("config is null")
)

// Commands 将子命令统一暴露给 main 包
var Commands = []*cli.Command{
	infoCmd,
	diagCmd,
}

var GlobalConfig *config.Config

var GlobalOptions = []cli.Flag{
	// set the global option by &cli.XXXFlag{}
}

var LoadGlobalOptions = func(ctx *cli.Context) error {
	// load the global option by ctx.XXX()
	return nil
}
