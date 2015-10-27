package mgo

import (
	. "github.com/eynstudio/gobreak/ddd"
	"github.com/eynstudio/gobreak/di"
)

func Init() {
	eventBus := NewEventBus()
	di.Root.MapAs(eventBus, (*EventBus)(nil))
	cmdBus := NewCmdBus()
	di.Root.MapAs(cmdBus, (*CmdBus)(nil))

	eventStore, _ := NewMongoEventStore()
	di.Root.ApplyAndMapAs(eventStore, (*EventStore)(nil))

	domainRepo := NewDomainRepo()
	di.Root.ApplyAndMapAs(domainRepo, (*DomainRepo)(nil))

	aggCmdHandler := NewAggCmdHandler()
	di.Root.ApplyAndMapAs(aggCmdHandler, (*AggCmdHandler)(nil))
	cmdBus.SetHandler(aggCmdHandler)
}
