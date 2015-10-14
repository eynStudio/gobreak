package main

import (
	"fmt"
	. "github.com/eynstudio/gobreak/ddd"
	"gopkg.in/mgo.v2/bson"
	"log"
)

func main() {
	eventBus := NewEventBus()
	eventBus.AddGlobalHandler(&LoggerSubscriber{})

	eventStore := NewMemoryEventStore(eventBus)

	repository, err := NewCallbackRepository(eventStore)
	if err != nil {
		log.Fatalf("could not create repository: %s", err)
	}

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

	invitationRepository := NewMemoryReadRepository()
	invitationProjector := NewInvitationProjector(invitationRepository)
	eventBus.AddHandler(invitationProjector, &InviteCreated{})
	eventBus.AddHandler(invitationProjector, &InviteAccepted{})
	eventBus.AddHandler(invitationProjector, &InviteDeclined{})

	eventID := NewGuid()
	guestListRepository := NewMemoryReadRepository()
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
	commandBus.HandleCmd(&CreateInvite{InvitationID: hadesID, Name: "Hades"})
	commandBus.HandleCmd(&AcceptInvite{InvitationID: hadesID})

	zeusID := NewGuid()
	commandBus.HandleCmd(&CreateInvite{InvitationID: zeusID, Name: "Zeus"})
	commandBus.HandleCmd(&DeclineInvite{InvitationID: zeusID})

	// Read all invites.
	invitations, _ := invitationRepository.FindAll()
	for _, i := range invitations {
		fmt.Printf("invitation: %#v\n", i)
	}

	// Read the guest list.
	guestList, _ := guestListRepository.Find(eventID)
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
