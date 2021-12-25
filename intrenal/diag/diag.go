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

	infoFilePathAbs, _ := filepath.Abs(conf.Diag.Input.DirPath)
	infoFileEntries, _ := os.ReadDir(infoFilePathAbs)
	for _, infoFileEntry := range infoFileEntries {
		infoFilePath := filepath.Join(infoFilePathAbs, infoFileEntry.Name())
		allInfos := ReadAllInfos(infoFilePath)
		for _, f := range allInfos {
			diagResult := diag.GetDiagResult(f)
			_logger.Info(diagResult)
		}
	}
}

// ReadAllInfos read infos from a node info file
func ReadAllInfos(infoFilePath string) []*info.AllInfo {
	infoFilePathAbs, _ := filepath.Abs(infoFilePath)

	// TODO add read info code
	log.Println(infoFilePathAbs)

	return nil
}
