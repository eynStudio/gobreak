package main

import (
	"fmt"
	. "github.com/eynstudio/gobreak/ddd"
)

type InvitationAggregate struct {
	*AggregateBase

	name     string
	age      int
	accepted bool
	declined bool
}

func (i *InvitationAggregate) AggType() string {
	return "Invitation"
}

func (i *InvitationAggregate) HandleCmd(command Cmd) error {
	switch command := command.(type) {
	case *CreateInvite:
		i.StoreEvent(&InviteCreated{command.InvitationID, command.Name, command.Age})
		return nil

	case *AcceptInvite:
		if i.name == "" {
			return fmt.Errorf("invitee does not exist")
		}

		if i.declined {
			return fmt.Errorf("%s already declined", i.name)
		}

		if i.accepted {
			return nil
		}

		i.StoreEvent(&InviteAccepted{i.ID()})
		return nil

	case *DeclineInvite:
		if i.name == "" {
			return fmt.Errorf("invitee does not exist")
		}

		if i.accepted {
			return fmt.Errorf("%s already accepted", i.name)
		}

		if i.declined {
			return nil
		}

		i.StoreEvent(&InviteDeclined{i.ID()})
		return nil
	}
	return fmt.Errorf("couldn't handle command")
}

func (i *InvitationAggregate) ApplyEvent(event Event) {
	switch event := event.(type) {
	case *InviteCreated:
		i.name = event.Name
		i.age = event.Age
	case *InviteAccepted:
		i.accepted = true
	case *InviteDeclined:
		i.declined = true
	}
}
