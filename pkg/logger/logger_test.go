package logger

import (
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"testing"
)

func TestFileLogger(t *testing.T) {

	o := config.OutputConfig{
		DirPath: "../../data/logger",
	}
	flog := GetFileLogger("test", o)
	flog.Info("hello")
}
