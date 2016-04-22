package cmdbus

import (
	"errors"
	. "github.com/eynstudio/gobreak"
)

var (
	ErrHandlerNotFound = errors.New("handler not found")
	handlers           = map[CmdHandler]bool{}
)

type Cmd interface {
	ID() GUID
}

type CmdHandler interface {
	CanHandle(cmd Cmd) bool
	Handle(Cmd) error
}

func SetHandler(handler CmdHandler) {
	handlers[handler] = true
}

func Exec(cmd Cmd) error {
	err := ErrHandlerNotFound
	for h := range handlers {
		if h.CanHandle(cmd) {
			err = h.Handle(cmd)
		}
	}
	return err
}
