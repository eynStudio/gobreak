package main

import (
	//	"fmt"
	. "github.com/eynstudio/gobreak/db/mgo"
	. "github.com/eynstudio/gobreak/ddd"
)

type Invitation struct {
	ID     GUID
	Name   string
	Status string
}

type InvitationProjector struct {
	repository MgoRepo
}

func NewInvitationProjector(repository MgoRepo) *InvitationProjector {
	p := &InvitationProjector{
		repository: repository,
	}
	return p
}

func (p *InvitationProjector) HandleEvent(event Event) {
	switch event := event.(type) {
	case *InviteCreated:
		i := &Invitation{
			ID:   event.InvitationID,
			Name: event.Name,
		}
		p.repository.Save(i.ID, i)
	case *InviteAccepted:
		m := p.repository.Get(event.InvitationID)
		i := m.(*Invitation)
		i.Status = "accepted"
		p.repository.Save(i.ID, i)
	case *InviteDeclined:
		m := p.repository.Get(event.InvitationID)
		i := m.(*Invitation)
		i.Status = "declined"
		p.repository.Save(i.ID, i)
	}
}

type GuestList struct {
	NumGuests   int
	NumAccepted int
	NumDeclined int
}

type GuestListProjector struct {
	repository MgoRepo
	eventID    GUID
}

func NewGuestListProjector(repository MgoRepo, eventID GUID) *GuestListProjector {
	p := &GuestListProjector{
		repository: repository,
		eventID:    eventID,
	}
	return p
}

func (p *GuestListProjector) HandleEvent(event Event) {
	switch event.(type) {
	case *InviteCreated:
		m := p.repository.Get(p.eventID)
		if m == nil {
			m = &GuestList{}
		}
		g := m.(*GuestList)
		p.repository.Save(p.eventID, g)
	case *InviteAccepted:
		m := p.repository.Get(p.eventID)
		g := m.(*GuestList)
		g.NumAccepted++
		p.repository.Save(p.eventID, g)
	case *InviteDeclined:
		m := p.repository.Get(p.eventID)
		g := m.(*GuestList)
		g.NumDeclined++
		p.repository.Save(p.eventID, g)
	}
}
