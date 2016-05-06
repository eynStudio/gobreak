package gobreak

import (
	"errors"
	"log"
)

type Error struct {
	Err  error
	done bool
}

func (p Error) IsErr() bool        { return p.Err != nil }
func (p Error) NotErr() bool       { return p.Err == nil }
func (p Error) LogErr()            { LogErr(p.Err) }
func (p *Error) ResetErr()         { p.Err = nil }
func (p *Error) SetErr(msg string) { p.Err = errors.New(msg) }

// SetDone if set done, NoErrExec not exec any more
func (p *Error) SetDone() { p.done = true }

//func (p Error) NoErrExec(f func()) { NoErrExec(p.Err, f) }

func (p *Error) SetErrIf(yes bool, msg string) {
	if yes {
		p.Err = errors.New(msg)

	}
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
