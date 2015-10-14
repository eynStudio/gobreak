package ddd

type EventHandler interface {
	HandleEvent(Event)
}

type EventBus interface {
	PublishEvent(Event)
	AddHandler(EventHandler, Event)
	AddGlobalHandler(EventHandler)
}

type eventBus struct {
	eventHandlers  map[string]map[EventHandler]bool
	globalHandlers map[EventHandler]bool
}

func NewEventBus() EventBus {
	return &eventBus{make(map[string]map[EventHandler]bool), make(map[EventHandler]bool)}
}

func (p *eventBus) PublishEvent(event Event) {
	if handlers, ok := p.eventHandlers[event.EventType()]; ok {
		for handler := range handlers {
			handler.HandleEvent(event)
		}
	}
	for handler := range p.globalHandlers {
		handler.HandleEvent(event)
	}
}

func (p *eventBus) AddHandler(handler EventHandler, event Event) {
	if _, ok := p.eventHandlers[event.EventType()]; !ok {
		p.eventHandlers[event.EventType()] = make(map[EventHandler]bool)
	}
	p.eventHandlers[event.EventType()][handler] = true
}

func (p *eventBus) AddGlobalHandler(handler EventHandler) {
	p.globalHandlers[handler] = true
}
