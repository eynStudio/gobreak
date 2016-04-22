package cmdbus

import (
	"errors"
	"log"

	. "github.com/eynstudio/gobreak/dddd/ddd"
)

var (
	ErrHandlerNotFound = errors.New("handler not found")
	handlers           = map[CmdHandler]bool{}
	_aggCmdHandler     = aggCmdHandler{}
)

func SetHandler(handler CmdHandler) {
	handlers[handler] = true
}

func Exec(cmd Cmd) error {
	log.Println("cmdbus.exec", cmd)

	if _aggCmdHandler.CanHandle(cmd) {
		return _aggCmdHandler.Handle(cmd)
	}

	log.Println("cmdbus.exec...........", cmd)

	err := ErrHandlerNotFound
	for h := range handlers {
		if h.CanHandle(cmd) {
			err = h.Handle(cmd)
		}
	}
	return err
}
