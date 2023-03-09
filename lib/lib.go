package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	json "github.com/json-iterator/go"
	"github.com/jwalton/gchalk"
	"github.com/lutfiharidha/logger-go/interfaces"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Loggers struct {
	namespace      string
	version        string
	printToConsole bool
	fileConfig     interfaces.LoggingFile
	level          interfaces.DebugLevel
	onLogger       func(msg interfaces.LoggerMessage, raw string)
	outputFormat   interfaces.OutputFormat
	color          *gchalk.Builder
	writer         io.Writer
}

func NewLib() interfaces.Logger {
	return &Loggers{
		printToConsole: true,
		outputFormat:   interfaces.OutputFormatDefault,
		level:          interfaces.DebugLevelVerbose,
	}
}

func (c *Loggers) New() interfaces.Logger {
	return NewLib()
}

func (c *Loggers) Init(namespace, version string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetOutput(c)
	c.color = gchalk.New()
	c.color.SetLevel(gchalk.LevelAnsi16m)
	if c.outputFormat == interfaces.OutputFormatJSON {
		c.color.SetLevel(gchalk.LevelNone)
	}
	c.namespace = namespace
	c.version = version
}

func (c *Loggers) ServiceName() string {
	return c.namespace
}

func (c *Loggers) ServiceVersion() string {
	return c.version
}

func (c *Loggers) DisableColor(val bool) {
	if val {
		c.color.SetLevel(gchalk.LevelNone)
	} else {
		c.color.SetLevel(gchalk.LevelAnsi16m)
	}

	if c.outputFormat == interfaces.OutputFormatJSON {
		c.color.SetLevel(gchalk.LevelNone)
	}
}

func (c *Loggers) SetLogFile(config interfaces.LoggingFile) {
	c.fileConfig = config
	if err := os.MkdirAll(filepath.Dir(c.fileConfig.Output), 0750); err != nil {
		c.Error(err)
	} else {
		if c.fileConfig.Enable {
			c.writer = &lumberjack.Logger{
				Filename: c.fileConfig.Output,
				MaxSize:  c.fileConfig.MaxSize,
				MaxAge:   c.fileConfig.MaxAge,
				Compress: c.fileConfig.Compress,
			}
		}

	}

}

func (c *Loggers) SetLogLevel(level interfaces.DebugLevel) {
	c.level = level
}

func (c *Loggers) GetLogLevel() (level interfaces.DebugLevel) {
	return c.level
}

func (c *Loggers) SetPrintToConsole(pr bool) {
	c.printToConsole = pr
}

func (c *Loggers) SetOnLoggerHandler(f func(msg interfaces.LoggerMessage, raw string)) {
	c.onLogger = f
}

func (c *Loggers) GetPrintToConsole() (pr bool) {
	return c.printToConsole
}

func (c *Loggers) GetOutputFormat() interfaces.OutputFormat {
	return c.outputFormat
}

func (c *Loggers) SetOutputFormat(op interfaces.OutputFormat) {
	c.outputFormat = op
}

func (c *Loggers) Trace(format interface{}, input ...interface{}) {
	c.output(c.createMsg(interfaces.LogLevelTrace, interfaces.GetCaller(2), format, input))
}

func (c *Loggers) Debug(format interface{}, input ...interface{}) {
	c.output(c.createMsg(interfaces.LogLevelDebug, interfaces.GetCaller(2), format, input))
}

func (c *Loggers) Notice(format interface{}, input ...interface{}) {
	c.output(c.createMsg(interfaces.LogLevelNotice, interfaces.GetCaller(2), format, input))
}

func (c *Loggers) Info(format interface{}, input ...interface{}) {
	c.output(c.createMsg(interfaces.LogLevelInfo, interfaces.GetCaller(2), format, input))
}

func (c *Loggers) Warning(format interface{}, input ...interface{}) {
	c.output(c.createMsg(interfaces.LogLevelWarning, interfaces.GetCaller(2), format, input))
}

func (c *Loggers) Success(format interface{}, input ...interface{}) {
	c.output(c.createMsg(interfaces.LogLevelSuccess, interfaces.GetCaller(2), format, input))
}

func (c *Loggers) Error(format interface{}, input ...interface{}) interfaces.Logger {
	c.output(c.createMsg(interfaces.LogLevelError, interfaces.GetCaller(2), format, input))
	return c
}

func (c *Loggers) Debugf(format string, input ...interface{}) {
	c.output(c.createMsg(interfaces.LogLevelTrace, interfaces.GetCaller(2), format, input))
}

func (c *Loggers) Errorf(format string, input ...interface{}) {
	c.output(c.createMsg(interfaces.LogLevelError, interfaces.GetCaller(2), format, input))
}

func (c *Loggers) Infof(format string, input ...interface{}) {
	c.output(c.createMsg(interfaces.LogLevelInfo, interfaces.GetCaller(2), format, input))
}

func (c *Loggers) Warnf(format string, input ...interface{}) {
	c.output(c.createMsg(interfaces.LogLevelWarning, interfaces.GetCaller(2), format, input))
}

func (c *Loggers) createMsg(level interfaces.LogLevel,
	caller interfaces.Caller,
	format interface{},
	input ...interface{}) (msg interfaces.LoggerMessage) {

	var inp []interface{}
	for _, s := range input {
		if val, ok := s.([]interface{}); ok {
			inp = append(inp, val...)
		}
	}

	var ffs string
	var msgS interface{}
	if val, ok := format.(string); ok {
		ffs = val
		msgS = fmt.Sprintf(ffs, inp...)
	} else if val, ok := format.(error); ok {
		ffs = val.Error()
		msgS = fmt.Sprintf(ffs, inp...)
	} else {
		msgS = format
	}

	return interfaces.LoggerMessage{
		ID:        uuid.New().String(),
		Time:      time.Now(),
		Level:     level,
		LevelName: interfaces.GetLogLevelString(level),
		File:      caller.File,
		Line:      caller.Line,
		FuncName:  caller.Fname,
		Message:   msgS,
	}
}

func (c *Loggers) createJSONOutput(msg interfaces.LoggerMessage) []byte {
	jsonOut := make(map[string]interface{})
	jsonOut["logid"] = msg.ID
	jsonOut["level"] = strings.ToLower(msg.LevelName)
	jsonOut["time"] = msg.Time.Format("2006-01-02T15:04:05.000-0700")

	jsonOut["caller"] = fmt.Sprintf("%s:%d", msg.File, msg.Line)

	vv := reflect.TypeOf(msg.Message)
	switch vv.Kind() {
	case reflect.String:
		jsonOut["message"] = msg.Message
	case reflect.Map:
		if val, ok := msg.Message.(map[string]interface{}); ok {
			for kk, vv := range val {
				jsonOut[kk] = vv
			}
		} else if val, ok := msg.Message.(map[string]string); ok {
			for kk, vv := range val {
				jsonOut[kk] = vv
			}
		} else {
			jsonOut["message"] = msg.Message
		}
	default:
		jsonOut["message"] = msg.Message
	}

	if val, err := json.Marshal(jsonOut); err == nil {
		return val
	}

	return nil
}

func (c *Loggers) store(msg interfaces.LoggerMessage, raw string) {
	if c.onLogger != nil {
		c.onLogger(msg, raw)
	}

	if c.outputFormat == interfaces.OutputFormatDefault {
		if c.printToConsole {
			fmt.Println(raw)
		}

		if c.fileConfig.Enable {
			if c.writer != nil {
				if !c.fileConfig.Json {
					ssd := raw + "\n"
					_, _ = c.writer.Write([]byte(ssd))
				} else {
					if out := c.createJSONOutput(msg); out != nil {
						ssd := string(out)
						if c.writer != nil {
							ssd := ssd + "\n"
							_, _ = c.writer.Write([]byte(ssd))
						}
					}
				}
			}
		}
	}

	if c.outputFormat == interfaces.OutputFormatJSON {
		if out := c.createJSONOutput(msg); out != nil {
			ssd := string(out)
			if c.printToConsole {
				fmt.Println(ssd)
			}
			if c.fileConfig.Enable {
				if c.writer != nil {
					ssd := ssd + "\n"
					_, _ = c.writer.Write([]byte(ssd))
				}
			}
		}
	}

}

func (c *Loggers) output(msg interfaces.LoggerMessage) {
	var raw string
	if c.outputFormat == interfaces.OutputFormatDefault {
		raw = c.ParsingLog(msg)
	}

	switch c.level {
	case interfaces.DebugLevelTrace:
		c.store(msg, raw)
	case interfaces.DebugLevelVerbose:
		switch msg.Level {
		case
			interfaces.LogLevelNotice,
			interfaces.LogLevelInfo,
			interfaces.LogLevelWarning,
			interfaces.LogLevelSuccess,
			interfaces.LogLevelError:
			c.store(msg, raw)
		}
	case interfaces.DebugLevelInfo:
		switch msg.Level {
		case interfaces.LogLevelNotice,
			interfaces.LogLevelInfo,
			interfaces.LogLevelWarning,
			interfaces.LogLevelSuccess, interfaces.LogLevelError:
			c.store(msg, raw)
		}
	case interfaces.DebugLevelWarning:
		switch msg.Level {
		case interfaces.LogLevelWarning, interfaces.LogLevelError:
			c.store(msg, raw)
		}
	case interfaces.DebugLevelError:
		switch msg.Level {
		case interfaces.LogLevelError:
			c.store(msg, raw)
		}
	}

}

func (c *Loggers) ParsingLog(msg interfaces.LoggerMessage) (raw string) {
	mm := c.color.WithBold()
	mmc := c.color.WithBold()
	var ems string
	var vms string
	mMsg := fmt.Sprintf("%s", msg.Message)

	vv := reflect.TypeOf(msg.Message)
	if vv != nil {
		switch vv.Kind() {
		case reflect.String:
		default:
			if values, err := json.Marshal(msg.Message); err == nil {
				mMsg = string(values)
			}
		}
	} else {
		if values, err := json.Marshal(msg.Message); err == nil {
			mMsg = string(values)
		}
	}

	switch msg.Level {
	case interfaces.LogLevelTrace:
		ems = mm.BrightWhite(interfaces.GetLogLevelPrintString(msg.Level))
		vms = mmc.White(mMsg)
	case interfaces.LogLevelDebug:
		ems = mm.BrightBlue(interfaces.GetLogLevelPrintString(msg.Level))
		vms = mmc.White(mMsg)
	case interfaces.LogLevelNotice:
		ems = mm.BrightCyan(interfaces.GetLogLevelPrintString(msg.Level))
		vms = mmc.BrightCyan(mMsg)
	case interfaces.LogLevelInfo:
		ems = mm.BrightMagenta(interfaces.GetLogLevelPrintString(msg.Level))
		vms = mmc.BrightMagenta(mMsg)
	case interfaces.LogLevelWarning:
		ems = mm.Yellow(interfaces.GetLogLevelPrintString(msg.Level))
		vms = mmc.Yellow(mMsg)
	case interfaces.LogLevelError:
		ems = mm.Red(interfaces.GetLogLevelPrintString(msg.Level))
		vms = mmc.Red(mMsg)
	case interfaces.LogLevelSuccess:
		ems = mm.Green(interfaces.GetLogLevelPrintString(msg.Level))
		vms = mmc.Green(mMsg)

	}

	raw = fmt.Sprintf("[%s][%s][%s][%s][%d] %s",
		c.color.Magenta(msg.Time.Format("3:04:05 PM")),
		ems, c.color.BrightWhite(c.namespace),
		c.color.BrightCyan(filepath.Base(msg.FuncName)),
		msg.Line,
		vms)
	return raw
}

func (c *Loggers) Quit() {
	os.Exit(0)
}

func (c *Loggers) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(p))
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) != 0 {
			c.Debug("%s", text)
		}
	}
	return len(p), nil
}

func (c *Loggers) NewSystemLogger() *log.Logger {
	logs := log.New(c, "", log.LstdFlags)
	logs.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	return logs
}

func (c *Loggers) Printf(f string, data ...interface{}) {
	c.Debug(f, data...)
}

func (c *Loggers) SetGinLog() gin.HandlerFunc {
	return func(ctxGin *gin.Context) {
		start := time.Now() // Start timer
		path := ctxGin.Request.URL.Path
		raw := ctxGin.Request.URL.RawQuery

		ctxGin.Next()

		param := gin.LogFormatterParams{}

		param.TimeStamp = time.Now() // Stop timer
		param.Latency = param.TimeStamp.Sub(start)
		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}
		param.ClientIP = ctxGin.ClientIP()
		param.Method = ctxGin.Request.Method
		param.StatusCode = ctxGin.Writer.Status()
		param.ErrorMessage = ctxGin.Errors.String()
		param.BodySize = ctxGin.Writer.Size()
		param.Path = ctxGin.FullPath()
		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path
		switch {
		case param.StatusCode >= 400 && param.StatusCode <= 499:
			{
				c.Warning("[GIN] %3d | %13v | %15s | %s %#v\n%s", param.StatusCode, param.Latency, param.ClientIP, param.Method, param.Path, param.ErrorMessage)
			}
		case param.StatusCode >= 500:
			{
				c.Error("[GIN] %3d | %13v | %15s | %s %#v\n%s", param.StatusCode, param.Latency, param.ClientIP, param.Method, param.Path, param.ErrorMessage)

			}
		default:
			c.Info("[GIN] %3d | %13v | %15s | %s %#v", param.StatusCode, param.Latency, param.ClientIP, param.Method, param.Path)
		}
	}
}
