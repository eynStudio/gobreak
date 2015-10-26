package ddd

import (
	"reflect"
)

type EventHandler interface {
	HandleEvent(Event)
}

type EventBus interface {
	PublishEvent(Event)
	AddHandler(EventHandler, ...Event)
	AddGlobalHandler(EventHandler)
}

type eventBus struct {
	eventHandlers  map[reflect.Type]map[EventHandler]bool
	globalHandlers map[EventHandler]bool
}

func NewEventBus() EventBus {
	return &eventBus{make(map[reflect.Type]map[EventHandler]bool),
		make(map[EventHandler]bool)}
}

func (p *eventBus) PublishEvent(event Event) {
	if handlers, ok := p.eventHandlers[reflect.TypeOf(event)]; ok {
		for handler := range handlers {
			handler.HandleEvent(event)
		}
	}

	for handler := range p.globalHandlers {
		handler.HandleEvent(event)
	}
}

func (p *eventBus) AddHandler(handler EventHandler, events ...Event) {
	for _, event := range events {
		evtType := reflect.TypeOf(event)
		if _, ok := p.eventHandlers[evtType]; !ok {
			p.eventHandlers[evtType] = make(map[EventHandler]bool)
		}
		p.eventHandlers[evtType][handler] = true
	}
}

func (p *eventBus) AddGlobalHandler(handler EventHandler) {
	p.globalHandlers[handler] = true
}
