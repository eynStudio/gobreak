package gobreak

import (
	"errors"
	"fmt"
	"log"
)

type Error struct {
	Err  error
	done bool
	msg  string
}

func (p Error) IsErr() bool                            { return p.Err != nil }
func (p Error) NotErr() bool                           { return p.Err == nil }
func (p Error) LogErr()                                { LogErr(p.Err) }
func (p *Error) ResetErr()                             { p.Err = nil }
func (p *Error) SetMsg(msg string)                     { p.msg = msg }
func (p *Error) SetMsgf(f string, args ...interface{}) { p.msg = fmt.Sprintf(f, args...) }
func (p *Error) SetErrf(f string, args ...interface{}) { p.SetErr(fmt.Sprintf(f, args...)) }
func (p *Error) SetErr(msg string) {
	p.SetMsg(msg)
	p.Err = errors.New(msg)
}

// SetDone if set done, NoErrExec not exec any more
func (p *Error) SetDone() { p.done = true }

func (p *Error) SetErrIf(yes bool, msg string) {
	if yes {
		p.SetErr(msg)
	}
}
func (p *Error) SetErrfIf(yes bool, f string, args ...interface{}) {
	if yes {
		p.SetErrf(f, args...)
	}
}

func (p Error) GetStatus() (m Status) {
	if p.IsErr() {
		m.ErrMsg(p.msg)
	} else {
		m.OkMsg(p.msg)
	}
	return
}

func (p Error) NoErrExec(f ...func()) {
	for _, it := range f {
		if !p.done {
			NoErrExec(p.Err, it)
		}
	}
}

func NoErrExec(err error, f func()) {
	if err == nil {
		f()
	}
}

func LogErr(err error) error {
	if err != nil {
		log.Printf("%#v", err)
	}
	return err
}
