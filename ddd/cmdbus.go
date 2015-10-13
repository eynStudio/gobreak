package ddd

import (
	"errors"
)

var (
	ErrHandlerExist    = errors.New("handler is already exist")
	ErrHandlerNotFound = errors.New("handler not found")
)

type CmdHandler interface {
	HandleCmd(Cmd) error
}

type CmdBus interface {
	HandleCmd(Cmd) error
	SetHandler(CmdHandler, Cmd) error
}

type cmdBus struct {
	handlers map[string]CmdHandler
}

func NewCmdBus() CmdBus {
	return &cmdBus{make(map[string]CmdHandler)}
}

func (p *cmdBus) HandleCmd(cmd Cmd) error {
	if handler, ok := p.handlers[cmd.CmdType()]; ok {
		return handler.HandleCmd(cmd)
	}
	return ErrHandlerNotFound
}

func (p *cmdBus) SetHandler(handler CmdHandler, cmd Cmd) error {
	if _, ok := p.handlers[cmd.CmdType()]; ok {
		return ErrHandlerExist
	}
	p.handlers[cmd.CmdType()] = handler
	return nil
}
