package mgo

import (
	"time"

	. "github.com/eynstudio/gobreak"
)

type LogInfo struct {
	Id    GUID      `Id`
	Time  time.Time `Time`
	Level int       `Level`
	Msg   string    `Msg`
}

type MgoLogger struct {
	repo MgoRepo
}

func (p *MgoLogger) Log(level int, msg string) {
	info := &LogInfo{
		Id:    p.repo.NewId(),
		Time:  time.Now(),
		Level: level,
		Msg:   msg,
	}
	p.repo.Save(info.Id, info)
}
