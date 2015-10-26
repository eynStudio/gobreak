package ddd

import (
	"errors"
)

var (
	ErrHandlerNotFound = errors.New("handler not found")
)

type CmdHandler interface {
	CanHandleCmd(cmd Cmd) bool
	HandleCmd(Cmd) error
}

type CmdBus interface {
	HandleCmd(Cmd) error
	SetHandler(CmdHandler)
}

type cmdBus struct {
	handlers map[CmdHandler]bool
}

func NewCmdBus() CmdBus {
	return &cmdBus{make(map[CmdHandler]bool)}
}

func (p *cmdBus) HandleCmd(cmd Cmd) error {
	err := ErrHandlerNotFound
	for handler := range p.handlers {
		if handler.CanHandleCmd(cmd) {
			err = handler.HandleCmd(cmd)
		}
	}
	return err
}

func (p *cmdBus) SetHandler(handler CmdHandler) {
	p.handlers[handler] = true
}
