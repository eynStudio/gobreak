package gobreak

import (
	"errors"
	"log"
)

type Error struct {
	Err error
}

func (p Error) IsErr() bool        { return p.Err != nil }
func (p Error) NotErr() bool       { return p.Err == nil }
func (p Error) LogErr()            { LogErr(p.Err) }
func (p Error) NoErrExec(f func()) { NoErrExec(p.Err, f) }
func (p *Error) ResetErr()         { p.Err = nil }
func (p *Error) SetErr(msg string) { p.Err = errors.New(msg) }
func (p *Error) SetErrIf(yes bool, msg string) {
	if yes {
		p.Err = errors.New(msg)

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
