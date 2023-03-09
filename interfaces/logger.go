package interfaces

import (
	"log"

	"github.com/gin-gonic/gin"
)

// Logger modules interface, using for dynamic modules
type Logger interface {
	New() Logger
	Init(namespace, version string)
	ServiceName() string
	ServiceVersion() string
	SetLogLevel(level DebugLevel)
	SetLogFile(LoggingFile)
	GetLogLevel() (level DebugLevel)
	SetPrintToConsole(pr bool)
	GetPrintToConsole() (pr bool)
	DisableColor(val bool)
	SetOnLoggerHandler(f func(msg LoggerMessage, raw string))
	SetOutputFormat(OutputFormat)
	GetOutputFormat() OutputFormat
	ParsingLog(msg LoggerMessage) (raw string)
	Write(p []byte) (int, error)
	Trace(format interface{}, input ...interface{})
	Debug(format interface{}, input ...interface{})
	Notice(format interface{}, input ...interface{})
	Info(format interface{}, input ...interface{})
	Warning(format interface{}, input ...interface{})
	Success(format interface{}, input ...interface{})
	Error(format interface{}, input ...interface{}) Logger
	Debugf(format string, input ...interface{})
	Errorf(format string, input ...interface{})
	Infof(format string, input ...interface{})
	Warnf(format string, input ...interface{})
	NewSystemLogger() *log.Logger
	Printf(string, ...interface{})
	Quit()
	SetGinLog() gin.HandlerFunc
}
