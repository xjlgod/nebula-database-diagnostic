package diag

import (
	"bufio"
	"encoding/json"
	"github.com/xjlgod/nebula-database-diagnostic/intrenal/info"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/diag"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/logger"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
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
			_logger.Infof("%s\n", strings.Join(diagResult, ""))
		}
	}
}

// ReadAllInfos read infos from a node info file
func ReadAllInfos(infoFilePath string) []*info.AllInfo {
	infoFilePathAbs, _ := filepath.Abs(infoFilePath)
	file, err := os.Open(infoFilePathAbs)
	defer file.Close()
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	rd := bufio.NewReader(file)
	allInfos := make([]*info.AllInfo, 0)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil {
			if io.EOF != err {
				log.Println(err.Error())
			}
			break
		}
		allInfo := new(info.AllInfo)
		err = json.Unmarshal([]byte(line), allInfo)
		if err != nil {
			log.Println(err.Error())
			break
		}
		allInfos = append(allInfos, allInfo)
	}

	return allInfos
}
