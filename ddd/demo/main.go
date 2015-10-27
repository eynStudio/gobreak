package main

import (
	"fmt"
	"log"
	//	"reflect"

	. "github.com/eynstudio/gobreak"
	. "github.com/eynstudio/gobreak/db/mgo"
	. "github.com/eynstudio/gobreak/ddd"
	. "github.com/eynstudio/gobreak/ddd/mgo"
	"github.com/eynstudio/gobreak/di"
	//	"gopkg.in/mgo.v2/bson"
)

type TestModel struct {
	Id   GUID
	Name string
}

func (p *TestModel) ID() GUID { return p.Id }

type TestCmdModel TestModel

func (p *TestCmdModel) ID() GUID { return p.Id }

type TestEventModel TestModel

func (p *TestEventModel) ID() GUID { return p.Id }
func main() {
	mgoCtx := NewMgoCtx(&MgoCfg{
		Server: "202.192.149.85:27017",
		Db:     "mis",
		User:   "mis",
		Pwd:    "mis.gs@stu.cn"})
	di.Root.Map(mgoCtx)

	Init()
	di.Root.Invoke(run)

}

func run(eventBus EventBus, commandBus CmdBus, domainRepo DomainRepo, aggCmdHandler AggCmdHandler) {
	eventBus.AddGlobalHandler(&LoggerSubscriber{})

	domainRepo.RegisterAggregate(&InvitationAggregate{}, NewInvitationAggregate)
	aggCmdHandler.SetAggregate(&InvitationAggregate{})

	invitationRepository := NewMgoRepo("test_Invitations", func() T { return &Invitation{} })
	di.Root.Apply(invitationRepository)
	invitationRepository.Clear()

	invitationProjector := NewInvitationProjector(invitationRepository)
	eventBus.AddHandler(invitationProjector)

	eventID := NewGuid()
	guestListRepository := NewMgoRepo("test_guest_lists", func() T { return &GuestList{} })
	di.Root.Apply(guestListRepository)
	guestListRepository.Clear()

	guestListProjector := NewGuestListProjector(guestListRepository, eventID)
	eventBus.AddHandler(guestListProjector)

	// Issue some invitations and responses.
	// Note that Athena tries to decline the event, but that is not allowed
	// by the domain logic in InvitationAggregate. The result is that she is
	// still accepted.
	athenaID := NewGuid()
	commandBus.HandleCmd(&CreateInvite{InvitationID: athenaID, Name: "Athena", Age: 42})
	commandBus.HandleCmd(&AcceptInvite{InvitationID: athenaID})
	err := commandBus.HandleCmd(&DeclineInvite{InvitationID: athenaID})
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	hadesID := NewGuid()
	commandBus.HandleCmd(&CreateInvite{InvitationID: hadesID, Name: "Hades", Age: 42})
	commandBus.HandleCmd(&AcceptInvite{InvitationID: hadesID})

	zeusID := NewGuid()
	commandBus.HandleCmd(&CreateInvite{InvitationID: zeusID, Name: "Zeus", Age: 42})
	commandBus.HandleCmd(&DeclineInvite{InvitationID: zeusID})

	testmodel := &TestModel{Id: zeusID, Name: "eyn"}
	commandBus.HandleCmd((*TestCmdModel)(testmodel))

	// Read all invites.
	invitations := invitationRepository.All()
	for _, i := range invitations {
		fmt.Printf("invitation: %#v\n", i)
	}

	// Read the guest list.
	guestList := guestListRepository.Get(eventID)
	fmt.Printf("guest list: %#v\n", guestList)
}

type LoggerSubscriber struct{}

func (l *LoggerSubscriber) HandleEvent(event Event) {
	log.Printf("event: %#v\n", event)
}
