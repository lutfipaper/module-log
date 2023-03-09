package interfaces

import (
	"encoding/json"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type LogLevel int

const (
	LogLevelTrace LogLevel = iota + 1
	LogLevelDebug
	LogLevelNotice
	LogLevelInfo
	LogLevelWarning
	LogLevelError
	LogLevelSuccess
)

type DebugLevel int

const (
	DebugLevelTrace DebugLevel = iota + 1
	DebugLevelVerbose
	DebugLevelInfo
	DebugLevelWarning
	DebugLevelError
)

func GetDebugLevelFromString(level string) DebugLevel {
	switch strings.ToLower(level) {
	case "trace":
		return DebugLevelTrace
	case "verbose":
		return DebugLevelVerbose
	case "info":
		return DebugLevelInfo
	case "warning":
		return DebugLevelWarning
	case "error":
		return DebugLevelError

	}

	return -1
}

var levelstring = map[LogLevel]string{
	LogLevelTrace:   "TRACE",
	LogLevelDebug:   "DEBUG",
	LogLevelNotice:  "NOTICE",
	LogLevelInfo:    "INFO",
	LogLevelWarning: "WARNING",
	LogLevelError:   "ERROR",
	LogLevelSuccess: "SUCCESS",
}

func GetLogLevelString(level LogLevel) string {
	return levelstring[level]
}

var levelPrintstring = map[LogLevel]string{
	LogLevelTrace:   "TRCE",
	LogLevelDebug:   "DBUG",
	LogLevelNotice:  "NTCE",
	LogLevelInfo:    "INFO",
	LogLevelWarning: "WARN",
	LogLevelError:   "EROR",
	LogLevelSuccess: "SUCS",
}

func GetLogLevelPrintString(level LogLevel) string {
	return levelPrintstring[level]
}

type LoggerMessage struct {
	ID        string      `json:"id"`
	Level     LogLevel    `json:"level"`
	LevelName string      `json:"levelName"`
	File      string      `json:"file"`
	Line      int         `json:"line"`
	FuncName  string      `json:"funcName"`
	Time      time.Time   `json:"time"`
	Message   interface{} `json:"message"`
}

type OutputFormat int

const (
	OutputFormatDefault OutputFormat = iota + 1
	OutputFormatJSON
)

func GetOutputFormatFromString(op string) OutputFormat {
	switch strings.ToLower(op) {
	case "default":
		return OutputFormatDefault
	case "json":
		return OutputFormatJSON
	}
	return OutputFormatDefault
}

type Caller struct {
	File       string `json:"file"`
	Line       int    `json:"line"`
	Fname      string `json:"fname"`
	FnameShort string `json:"-"`
}

func (c Caller) String() string {
	bs, _ := json.Marshal(c)
	return string(bs)
}

func GetCaller(skip int) (cs Caller) {
	if skip == 0 {
		skip = 1
	}

	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		cs.File = file
		cs.Line = line
		fc := runtime.FuncForPC(pc)
		cs.Fname = fc.Name()
		drname := filepath.Dir(fc.Name())
		spsd := strings.Split(drname, "/")
		if len(spsd) >= 2 {
			spsd = spsd[len(spsd)-2:]
		}
		spsd = append(spsd, filepath.Base(fc.Name()))
		cs.FnameShort = filepath.Join(spsd...)
	}

	return cs

}
