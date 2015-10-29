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

type Log interface {
	Trace(msg string)
	Tracef(format string, v ...interface{})
	Debug(msg string)
	Debugf(format string, v ...interface{})
	Info(msg string)
	Infof(format string, v ...interface{})
	Warn(msg string)
	Warnf(format string, v ...interface{})
	Error(msg string)
	Errorf(format string, v ...interface{})
	Fatal(msg string)
	Fatalf(format string, v ...interface{})
}

type Logs struct {
}

type log struct {
}

func NewLog() Log { return &log{} }

func (p *log) Log(level int, msg string) { fmt.Printf("%d:%s\n", level, msg) }
func (p *log) Logf(level int, format string, v ...interface{}) {
	p.Log(level, fmt.Sprintf(format, v...))
}

func (p *log) Trace(msg string)                       { p.Log(LevelTrace, msg) }
func (p *log) Tracef(format string, v ...interface{}) { p.Logf(LevelTrace, format, v...) }
