package main

import (
	"fmt"
	"log"

	. "github.com/eynstudio/gobreak"
	. "github.com/eynstudio/gobreak/db/mgo"
	. "github.com/eynstudio/gobreak/ddd"
	. "github.com/eynstudio/gobreak/ddd/mgo"
	"github.com/eynstudio/gobreak/di"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	mgoCtx := NewMgoCtx(&MgoCfg{
		Server: "202.192.149.85:27017",
		Db:     "mis",
		User:   "mis",
		Pwd:    "mis.gs@stu.cn"})
	di.Root.Map(mgoCtx)

	eventRepo := NewMgoRepo("DomainEvents", NewMongoAggregateRecord)
	di.Root.Apply(eventRepo)

	eventRepo.Clear()

	eventBus := NewEventBus()
	eventBus.AddGlobalHandler(&LoggerSubscriber{})

	eventStore, _ := NewMongoEventStore(eventBus, eventRepo)
	di.Root.MapAs(eventStore, (*EventStore)(nil))

	repository := NewDomainRepo()
	di.Root.Apply(repository)

	eventStore.RegisterEventType(&InviteCreated{}, func() Event { return &InviteCreated{} })
	eventStore.RegisterEventType(&InviteAccepted{}, func() Event { return &InviteAccepted{} })
	eventStore.RegisterEventType(&InviteDeclined{}, func() Event { return &InviteDeclined{} })

	repository.RegisterAggregate(&InvitationAggregate{},
		func(id GUID) Aggregate {
			return &InvitationAggregate{
				AggregateBase: NewAggregateBase(id),
			}
		},
	)

	handler, err := NewAggregateCommandHandler(repository)
	if err != nil {
		log.Fatalf("could not create command handler: %s", err)
	}

	handler.SetAggregate(&InvitationAggregate{}, &CreateInvite{})
	handler.SetAggregate(&InvitationAggregate{}, &AcceptInvite{})
	handler.SetAggregate(&InvitationAggregate{}, &DeclineInvite{})

	commandBus := NewCmdBus()
	commandBus.SetHandler(handler, &CreateInvite{})
	commandBus.SetHandler(handler, &AcceptInvite{})
	commandBus.SetHandler(handler, &DeclineInvite{})

	invitationRepository := NewMgoRepo("test_Invitations", func() T { return &Invitation{} })
	di.Root.Apply(invitationRepository)
	invitationRepository.Clear()

	invitationProjector := NewInvitationProjector(invitationRepository)
	eventBus.AddHandler(invitationProjector, &InviteCreated{})
	eventBus.AddHandler(invitationProjector, &InviteAccepted{})
	eventBus.AddHandler(invitationProjector, &InviteDeclined{})

	eventID := NewGuid()
	guestListRepository := NewMgoRepo("test_guest_lists", func() T { return &GuestList{} })
	di.Root.Apply(guestListRepository)
	guestListRepository.Clear()

	guestListProjector := NewGuestListProjector(guestListRepository, eventID)
	eventBus.AddHandler(guestListProjector, &InviteCreated{})
	eventBus.AddHandler(guestListProjector, &InviteAccepted{})
	eventBus.AddHandler(guestListProjector, &InviteDeclined{})

	// Issue some invitations and responses.
	// Note that Athena tries to decline the event, but that is not allowed
	// by the domain logic in InvitationAggregate. The result is that she is
	// still accepted.
	athenaID := NewGuid()
	commandBus.HandleCmd(&CreateInvite{InvitationID: athenaID, Name: "Athena", Age: 42})
	commandBus.HandleCmd(&AcceptInvite{InvitationID: athenaID})
	err = commandBus.HandleCmd(&DeclineInvite{InvitationID: athenaID})
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	hadesID := NewGuid()
	commandBus.HandleCmd(&CreateInvite{InvitationID: hadesID, Name: "Hades", Age: 42})
	commandBus.HandleCmd(&AcceptInvite{InvitationID: hadesID})

	zeusID := NewGuid()
	commandBus.HandleCmd(&CreateInvite{InvitationID: zeusID, Name: "Zeus", Age: 42})
	commandBus.HandleCmd(&DeclineInvite{InvitationID: zeusID})

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

func NewGuid() GUID {
	id := bson.NewObjectId().Hex()
	return GUID(id)
}
