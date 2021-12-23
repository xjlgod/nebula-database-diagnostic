package logger

import "fmt"

var (
	// logger.CMDLogger.Info()
	CMDLogger  Logger = DefaultLogger{false}
	FileLogger Logger = DefaultLogger{true}
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
	logToFile bool
}

func (d DefaultLogger) Info(msg ...interface{}) {
	d.info(fmt.Sprint(msg))
}
func (d DefaultLogger) Infof(format string, msg ...interface{}) {
	d.info(fmt.Sprintf(format, msg))
}
func (d DefaultLogger) info(msg string) {
	if d.logToFile {
		// log to file
	} else {
		//log to cmd
	}
}
func (d DefaultLogger) Warn(msg ...interface{}) {
	d.warn(fmt.Sprint(msg))
}
func (d DefaultLogger) Warnf(format string, msg ...interface{}) {
	d.warn(fmt.Sprintf(format, msg))
}
func (d DefaultLogger) warn(msg string) {
	if d.logToFile {
		// log to file
	} else {
		//log to cmd
	}
}
func (d DefaultLogger) Error(msg ...interface{}) {
	d.error(fmt.Sprint(msg))
}
func (d DefaultLogger) Errorf(format string, msg ...interface{}) {
	d.error(fmt.Sprintf(format, msg))
}
func (d DefaultLogger) error(msg string) {
	if d.logToFile {
		// log to file
	} else {
		//log to cmd
	}
}
func (d DefaultLogger) Fatal(msg ...interface{}) {
	d.fatal(fmt.Sprint(msg))
}
func (d DefaultLogger) Fatalf(format string, msg ...interface{}) {
	d.fatal(fmt.Sprintf(format, msg))
}
func (d DefaultLogger) fatal(msg string) {
	if d.logToFile {
		// log to file
	} else {
		//log to cmd
	}
}
