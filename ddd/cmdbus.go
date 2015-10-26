package ddd

import (
	"errors"
	"reflect"
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
	handlers map[reflect.Type]CmdHandler
}

func NewCmdBus() CmdBus {
	return &cmdBus{make(map[reflect.Type]CmdHandler)}
}

func (p *cmdBus) HandleCmd(cmd Cmd) error {
	err := ErrHandlerNotFound
	for _, handler := range p.handlers {
		if handler.CanHandleCmd(cmd) {
			err = handler.HandleCmd(cmd)
		}
	}
	return err
}

func (p *cmdBus) SetHandler(handler CmdHandler) {
	p.handlers[reflect.TypeOf(handler)] = handler
}
