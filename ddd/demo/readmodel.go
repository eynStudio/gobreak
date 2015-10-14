package main

import (
	. "github.com/eynstudio/gobreak/ddd"
)

type Invitation struct {
	ID     GUID
	Name   string
	Status string
}

type InvitationProjector struct {
	repository ReadRepository
}

func NewInvitationProjector(repository ReadRepository) *InvitationProjector {
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
		m, _ := p.repository.Find(event.InvitationID)
		i := m.(*Invitation)
		i.Status = "accepted"
		p.repository.Save(i.ID, i)
	case *InviteDeclined:
		m, _ := p.repository.Find(event.InvitationID)
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
	repository ReadRepository
	eventID    GUID
}

func NewGuestListProjector(repository ReadRepository, eventID GUID) *GuestListProjector {
	p := &GuestListProjector{
		repository: repository,
		eventID:    eventID,
	}
	return p
}

func (p *GuestListProjector) HandleEvent(event Event) {
	switch event.(type) {
	case *InviteCreated:
		m, _ := p.repository.Find(p.eventID)
		if m == nil {
			m = &GuestList{}
		}
		g := m.(*GuestList)
		p.repository.Save(p.eventID, g)
	case *InviteAccepted:
		m, _ := p.repository.Find(p.eventID)
		g := m.(*GuestList)
		g.NumAccepted++
		p.repository.Save(p.eventID, g)
	case *InviteDeclined:
		m, _ := p.repository.Find(p.eventID)
		g := m.(*GuestList)
		g.NumDeclined++
		p.repository.Save(p.eventID, g)
	}
}
