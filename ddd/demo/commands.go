package main

import (
	. "github.com/eynstudio/gobreak"
)

type CreateInvite struct {
	InvitationID GUID
	Name         string
	Age          int
}

func (c *CreateInvite) ID() GUID { return c.InvitationID }

type AcceptInvite struct {
	InvitationID GUID
}

func (c *AcceptInvite) ID() GUID { return c.InvitationID }

type DeclineInvite struct {
	InvitationID GUID
}

func (c *DeclineInvite) ID() GUID { return c.InvitationID }
