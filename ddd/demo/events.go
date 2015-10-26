package main

import (
	. "github.com/eynstudio/gobreak"
)

type InviteCreated struct {
	InvitationID GUID
	Name         string
	Age          int
}

func (c *InviteCreated) ID() GUID { return c.InvitationID }

type InviteAccepted struct {
	InvitationID GUID
}

func (c *InviteAccepted) ID() GUID { return c.InvitationID }

type InviteDeclined struct {
	InvitationID GUID
}

func (c *InviteDeclined) ID() GUID { return c.InvitationID }
