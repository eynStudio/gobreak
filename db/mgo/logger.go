package mgo

import (
	"time"

	. "github.com/eynstudio/gobreak"
	"github.com/eynstudio/gobreak/log"
)

type LogInfo struct {
	Id     GUID      `Id`
	Time   time.Time `Time`
	Level  string    `Level`
	Msg    string    `Msg`
}

type MgoLogger struct {
	repo MgoRepo
}

func NewMgoLogger(repo MgoRepo) *MgoLogger {
	return &MgoLogger{repo}
}

func (p *MgoLogger) Log(level log.LogLevel, msg string) {
	info := &LogInfo{
		Id:    p.repo.NewId(),
		Time:  time.Now(),
		Level: log.LogLevelName[level],
		Msg:   msg,
	}
	p.repo.Save(info.Id, info)
}
