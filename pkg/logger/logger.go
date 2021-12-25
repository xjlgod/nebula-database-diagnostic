package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xjlgod/nebula-database-diagnostic/pkg/config"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var (
	cmdLoggers  = make(map[string]*DefaultLogger, 0)
	cmdMux      sync.RWMutex
	fileLoggers = make(map[string]*DefaultLogger, 0)
	fileMux     sync.RWMutex
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
		Filepath() string
		Filename() string
	}
)

// DefaultLogger impl the default logger by logrus
type DefaultLogger struct {
	logr *logrus.Logger

	logToFile bool
	filepath  string
	filename string
}

func (d *DefaultLogger) Filepath() string {
	return d.filepath
}

func (d *DefaultLogger) Filename() string {
	return d.filename
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

func GetCmdLogger(n string) Logger {
	cmdMux.Lock()
	if _, ok := cmdLoggers[n]; !ok {
		initCmdLogger(n)
	}
	cmdMux.Unlock()

	cmdMux.RLock()
	defer cmdMux.RUnlock()
	return cmdLoggers[n]
}

func initCmdLogger(n string) {
	logr := logrus.New()
	logr.SetFormatter(&logrus.TextFormatter{})
	logr.SetOutput(os.Stdout)
	logr.SetLevel(logrus.InfoLevel)

	cmdLogger := new(DefaultLogger)
	cmdLogger.logToFile = false
	cmdLogger.logr = logr
	cmdLoggers[n] = cmdLogger

	timeUnix := time.Now().Unix()
	filename := fmt.Sprintf("%s_%s", n, strconv.FormatInt(timeUnix, 10))
	cmdLogger.filename = filename

}

func GetFileLogger(n string, o config.OutputConfig) Logger {
	fileMux.Lock()
	if _, ok := fileLoggers[n]; !ok {
		initFileLogger(n, o)
	}
	fileMux.Unlock()

	fileMux.RLock()
	defer fileMux.RUnlock()
	return fileLoggers[n]
}

func initFileLogger(n string, o config.OutputConfig) {
	logr := logrus.New()
	logr.SetFormatter(&logrus.TextFormatter{})
	timeUnix := time.Now().Unix()
	p, _ := filepath.Abs(o.DirPath)
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		os.Mkdir(p, os.ModePerm)
	}
	fileLogger := new(DefaultLogger)
	fileLogger.logToFile = true

	filename := fmt.Sprintf("%s_%s", n, strconv.FormatInt(timeUnix, 10))
	fileLogger.filename = filename
	fileLogger.filepath = filepath.Join(p, filename + ".log")
	file, err := os.OpenFile(fileLogger.filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logr.Fatal(err)
	}
	writer := io.Writer(file)
	logr.SetOutput(writer)
	logr.SetLevel(logrus.InfoLevel)

	fileLogger.logr = logr
	fileLoggers[n] = fileLogger
}
