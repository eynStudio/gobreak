package main

import (
	"fmt"

	. "github.com/eynstudio/gobreak/ddd"
		. "github.com/eynstudio/gobreak"
)

const (
	InvitationAggType = "Invitation"
)

type InvitationState struct {
	name     string
	age      int
	accepted bool
	declined bool
}

type InvitationAggregate struct {
	*AggregateBase
	StateModel InvitationState
}

func (i *InvitationAggregate) AggType() string {
	return InvitationAggType
}

func (i *InvitationAggregate) HandleCmd(command Cmd) error {

	switch command := command.(type) {
	case *CreateInvite:
		i.StoreEvent(i.ApplyEvent(&InviteCreated{command.InvitationID, command.Name, command.Age}))
		return nil

	case *AcceptInvite:
		if i.StateModel.name == "" {
			return fmt.Errorf("invitee does not exist")
		}

		if i.StateModel.declined {
			return fmt.Errorf("%s already declined", i.StateModel.name)
		}

		if i.StateModel.accepted {
			return nil
		}

		i.StoreEvent(i.ApplyEvent(&InviteAccepted{i.ID()}))
		return nil

	case *DeclineInvite:
		if i.StateModel.name == "" {
			return fmt.Errorf("invitee does not exist")
		}

		if i.StateModel.accepted {
			return fmt.Errorf("%s already accepted", i.StateModel.name)
		}

		if i.StateModel.declined {
			return nil
		}

		i.StoreEvent(i.ApplyEvent( &InviteDeclined{i.ID()}))
		return nil
	}
	return fmt.Errorf("couldn't handle command")
}

func (i *InvitationAggregate) ApplyEvent(event Event) Event {
	fmt.Println("ApplyEvent",event)
	switch event := event.(type) {
	case *InviteCreated:
		i.StateModel.name = event.Name
		i.StateModel.age = event.Age
	case *InviteAccepted:
		i.StateModel.accepted = true
	case *InviteDeclined:
		i.StateModel.declined = true
	}
	return event
}

func (i *InvitationAggregate) 	GetSnapshot() T{
	return &i.StateModel
}