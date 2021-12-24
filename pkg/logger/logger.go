package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	// logger.CMDLogger.Info()
	CMDLogger  *DefaultLogger = &DefaultLogger{logToFile: false}
	FileLogger *DefaultLogger = &DefaultLogger{logToFile: true}
)

type (
	Logger interface {
		Info(...interface{})
		Infof(string, ...interface{})
		Warn(...interface{})
		Warnf(string, ...interface{})
		Error(...interface{})
		Errorf(string, ...interface{})
		Fatal(...interface{})
		Fatalf(string, ...interface{})
	}
)

// DefaultLogger impl the default logger by logrus
type DefaultLogger struct {
	logr *logrus.Logger

	logToFile bool
	filepath  string
}

func (d *DefaultLogger) Info(msg ...interface{}) {
	d.info(fmt.Sprint(msg...))
}
func (d *DefaultLogger) Infof(format string, msg ...interface{}) {
	d.info(fmt.Sprintf(format, msg...))
}
func (d *DefaultLogger) info(msg string) {
	d.logr.Info(msg)
}
func (d *DefaultLogger) Warn(msg ...interface{}) {
	d.warn(fmt.Sprint(msg...))
}
func (d *DefaultLogger) Warnf(format string, msg ...interface{}) {
	d.warn(fmt.Sprintf(format, msg...))
}

func (d *DefaultLogger) warn(msg string) {
	d.logr.Warn(msg)
}
func (d *DefaultLogger) Error(msg ...interface{}) {
	d.error(fmt.Sprint(msg...))
}
func (d *DefaultLogger) Errorf(format string, msg ...interface{}) {
	d.error(fmt.Sprintf(format, msg...))
}
func (d *DefaultLogger) error(msg string) {
	d.logr.Error(msg)
}
func (d *DefaultLogger) Fatal(msg ...interface{}) {
	d.fatal(fmt.Sprint(msg...))
}
func (d *DefaultLogger) Fatalf(format string, msg ...interface{}) {
	d.fatal(fmt.Sprintf(format, msg...))
}
func (d *DefaultLogger) fatal(msg string) {
	d.logr.Fatal(msg)
}

func InitCmdLogger() {

	logr := logrus.New()
	logr.SetFormatter(&logrus.TextFormatter{})
	logr.SetOutput(os.Stdout)
	logr.SetLevel(logrus.InfoLevel)

	CMDLogger.logr = logr

}

func InitFileLogger(o config.OutputConfig) {

	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{})
	timeUnix := time.Now().Unix()
	FileLogger.filepath = filepath.Join(o.DirPath, strconv.FormatInt(timeUnix, 10)+".log")
	file, err := os.OpenFile(FileLogger.filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	writer := io.Writer(file)
	log.SetOutput(writer)
	log.SetLevel(logrus.InfoLevel)
	FileLogger.logr = log

}
