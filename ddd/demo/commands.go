package main

import (
	. "github.com/eynstudio/gobreak/ddd"
)

type CreateInvite struct {
	InvitationID GUID
	Name         string
	Age          int
}

func (c *CreateInvite) ID() GUID        { return c.InvitationID }
func (c *CreateInvite) AggType() string { return "Invitation" }
func (c *CreateInvite) CmdType() string { return "CreateInvite" }

type AcceptInvite struct {
	InvitationID GUID
}

func (c *AcceptInvite) ID() GUID        { return c.InvitationID }
func (c *AcceptInvite) AggType() string { return "Invitation" }
func (c *AcceptInvite) CmdType() string { return "AcceptInvite" }

type DeclineInvite struct {
	InvitationID GUID
}

func (c *DeclineInvite) ID() GUID        { return c.InvitationID }
func (c *DeclineInvite) AggType() string { return "Invitation" }
func (c *DeclineInvite) CmdType() string { return "DeclineInvite" }
