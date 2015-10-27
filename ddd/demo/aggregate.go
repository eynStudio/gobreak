package main

import (
	"fmt"

	. "github.com/eynstudio/gobreak"
	. "github.com/eynstudio/gobreak/ddd"
)

type InvitationState struct {
	Name     string `Name`
	Age      int    `Age`
	Accepted bool   `Accepted`
	Declined bool   `Edclined`
}

type InvitationAggregate struct {
	*AggregateBase
	StateModel InvitationState
}

func NewInvitationAggregate(id GUID) Aggregate {
	return &InvitationAggregate{
		AggregateBase: NewAggregateBase(id),
	}
}

func (i *InvitationAggregate) RegistedCmds() []Cmd {
	return []Cmd{&CreateInvite{}, &AcceptInvite{}, &DeclineInvite{}, &TestCmdModel{}}
}

func (i *InvitationAggregate) HandleCmd(command Cmd) error {

	switch command := command.(type) {
	case *TestCmdModel:
		fmt.Println("TestCmdModel")
		fmt.Println(command)
		i.ApplyEvent((*TestEventModel)(command))
		return nil
	case *CreateInvite:
		i.ApplyEvent(&InviteCreated{command.InvitationID, command.Name, command.Age})
		return nil

	case *AcceptInvite:
		if i.StateModel.Name == "" {
			return fmt.Errorf("invitee does not exist")
		}

		if i.StateModel.Declined {
			return fmt.Errorf("%s already declined", i.StateModel.Name)
		}

		if i.StateModel.Accepted {
			return nil
		}

		i.ApplyEvent(&InviteAccepted{i.ID()})
		return nil

	case *DeclineInvite:
		if i.StateModel.Name == "" {
			return fmt.Errorf("invitee does not exist")
		}

		if i.StateModel.Accepted {
			return fmt.Errorf("%s already accepted", i.StateModel.Name)
		}

		if i.StateModel.Declined {
			return nil
		}

		i.ApplyEvent(&InviteDeclined{i.ID()})
		return nil
	}
	return fmt.Errorf("couldn't handle command")
}

func (i *InvitationAggregate) ApplyEvent(event Event) {
	switch evt := event.(type) {
	case *TestEventModel:
		fmt.Println("TestEventModel")
		fmt.Println(evt)
	case *InviteCreated:
		i.StateModel.Name = evt.Name
		i.StateModel.Age = evt.Age
	case *InviteAccepted:
		i.StateModel.Accepted = true
	case *InviteDeclined:
		i.StateModel.Declined = true
	}
	i.IncrementVersion()
	i.StoreEvent(event)
}

func (i *InvitationAggregate) GetSnapshot() T {
	return &i.StateModel
}
