package log

import (
	"fmt"
)

const (
	LevelFatal = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

var loggers []Logger

type Logger interface {
	Log(level int, msg string)
}

func AddLogger(l Logger) { loggers = append(loggers, l) }

func Log(level int, msg string) {
	for _, it := range loggers {
		it.Log(level, msg)
	}
}

func Logf(level int, format string, v ...interface{}) { Log(level, fmt.Sprintf(format, v...)) }
func Trace(msg string)                                { Log(LevelTrace, msg) }
func Tracef(format string, v ...interface{})          { Logf(LevelTrace, format, v...) }
func Debug(msg string)                                { Log(LevelDebug, msg) }
func Debugf(format string, v ...interface{})          { Logf(LevelDebug, format, v...) }
func Info(msg string)                                 { Log(LevelInfo, msg) }
func Infof(format string, v ...interface{})           { Logf(LevelInfo, format, v...) }
func Warn(msg string)                                 { Log(LevelWarn, msg) }
func Warnf(format string, v ...interface{})           { Logf(LevelWarn, format, v...) }
func Error(msg string)                                { Log(LevelError, msg) }
func Errorf(format string, v ...interface{})          { Logf(LevelError, format, v...) }
func Fatal(msg string)                                { Log(LevelFatal, msg) }
func Fatalf(format string, v ...interface{})          { Logf(LevelFatal, format, v...) }
