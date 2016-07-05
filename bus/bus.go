package bus

import (
	"fmt"
	"reflect"
	"sync"
)

var _bus *Bus = New()

type Bus struct {
	handlers map[string][]*handler
	lock     sync.Mutex
	wg       sync.WaitGroup
}

type handler struct {
	callBack reflect.Value
	once     bool
	async    bool
	serial   bool //run serially (true) or concurrently (false)
	sync.Mutex
}

func New() *Bus { return &Bus{make(map[string][]*handler), sync.Mutex{}, sync.WaitGroup{}} }
func newH(fn interface{}, once, async, serial bool) *handler {
	return &handler{reflect.ValueOf(fn), once, async, serial, sync.Mutex{}}
}

func Sub(topic string, fn interface{}) error     { return _bus.Sub(topic, fn) }
func SubOnce(topic string, fn interface{}) error { return _bus.SubOnce(topic, fn) }
func SubAsync(topic string, fn interface{}, serial bool) error {
	return _bus.SubAsync(topic, fn, serial)
}
func SubOnceAsync(topic string, fn interface{}) error { return _bus.SubOnceAsync(topic, fn) }
func UnSub(topic string, handler interface{})         { _bus.UnSub(topic, handler) }
func Pub(topic string, args ...interface{})           { _bus.Pub(topic, args...) }
func WaitAsync()                                      { _bus.wg.Wait() }

func (p *Bus) Sub(topic string, fn interface{}) error {
	return p.sub(topic, fn, newH(fn, false, false, false))
}

func (p *Bus) SubAsync(topic string, fn interface{}, serial bool) error {
	return p.sub(topic, fn, newH(fn, false, true, serial))
}

func (p *Bus) SubOnce(topic string, fn interface{}) error {
	return p.sub(topic, fn, newH(fn, true, false, false))
}

func (p *Bus) SubOnceAsync(topic string, fn interface{}) error {
	return p.sub(topic, fn, newH(fn, true, true, false))
}

func (p *Bus) sub(topic string, fn interface{}, h *handler) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if !(reflect.TypeOf(fn).Kind() == reflect.Func) {
		return fmt.Errorf("%s is not reflect.Func", reflect.TypeOf(fn).Kind())
	}
	p.handlers[topic] = append(p.handlers[topic], h)
	return nil
}

func (p *Bus) UnSub(topic string, handler interface{}) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if _, ok := p.handlers[topic]; ok && len(p.handlers[topic]) > 0 {
		p.remHandler(topic, reflect.ValueOf(handler))
	}
}

func (p *Bus) Pub(topic string, args ...interface{}) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if handlers, ok := p.handlers[topic]; ok {
		for _, handler := range handlers {
			if !handler.async {
				p.pub(handler, topic, args...)
			} else {
				p.wg.Add(1)
				go p.pubAsync(handler, topic, args...)
			}
		}
	}
}

func (p *Bus) pub(handler *handler, topic string, args ...interface{}) {
	passedArguments := make([]reflect.Value, 0)
	for _, arg := range args {
		passedArguments = append(passedArguments, reflect.ValueOf(arg))
	}
	handler.callBack.Call(passedArguments)

	if handler.once {
		p.remHandler(topic, handler.callBack)
	}
}

func (p *Bus) pubAsync(handler *handler, topic string, args ...interface{}) {
	defer p.wg.Done()
	if handler.serial {
		handler.Lock()
		defer handler.Unlock()
	}
	p.pub(handler, topic, args...)
}

func (p *Bus) remHandler(topic string, callback reflect.Value) {
	i := p.findHandlerIdx(topic, callback)
	if i >= 0 {
		p.handlers[topic] = append(p.handlers[topic][:i], p.handlers[topic][i+1:]...)
	}
}

func (p *Bus) findHandlerIdx(topic string, callback reflect.Value) int {
	if _, ok := p.handlers[topic]; ok {
		for idx, handler := range p.handlers[topic] {
			if handler.callBack == callback {
				return idx
			}
		}
	}
	return -1
}

func (p *Bus) WaitAsync() { p.wg.Wait() }
