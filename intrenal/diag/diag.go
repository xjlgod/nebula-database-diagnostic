package diag

import (
	"github.com/xjlgod/nebula-database-diagnostic/intrenal/info"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/diag"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/logger"
	"log"
	"os"
	"path/filepath"
)

func Run(conf *config.Config) {
	var _logger logger.Logger
	if conf.Diag.Output.LogToFile {
		_logger = logger.GetFileLogger("diag", conf.Diag.Output)
	} else {
		_logger = logger.GetCmdLogger("diag")
	}

	allInfos := ReadAllInfos(conf)
	for _, f := range allInfos {
		diagResult := diag.GetDiagResult(f)
		_logger.Info(diagResult)
	}
}

// ReadAllInfos read infos from a dir which included nodes info data
func ReadAllInfos(conf *config.Config) []*info.AllInfo {
	inputDirPath := conf.Diag.Input.DirPath
	inputDirPathAbs, _ := filepath.Abs(inputDirPath)
	infoFileEntries, _ := os.ReadDir(inputDirPathAbs)

	for _, entry := range infoFileEntries {
		log.Println(entry.Name())
		// TODO add read
	}

	return nil
}
