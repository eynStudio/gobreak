package main

import (
	. "github.com/eynstudio/gobreak/ddd"
)

type InviteCreated struct {
	InvitationID GUID   `bson:"_d"`
	Name         string `Name`
	Age          int    `Age`
}

func (c *InviteCreated) ID() GUID          { return c.InvitationID }
func (c *InviteCreated) AggType() string   { return "Invitation" }
func (c *InviteCreated) EventType() string { return "InviteCreated" }

type InviteAccepted struct {
	InvitationID GUID
}

func (c *InviteAccepted) ID() GUID          { return c.InvitationID }
func (c *InviteAccepted) AggType() string   { return "Invitation" }
func (c *InviteAccepted) EventType() string { return "InviteAccepted" }

type InviteDeclined struct {
	InvitationID GUID
}

func (c *InviteDeclined) ID() GUID          { return c.InvitationID }
func (c *InviteDeclined) AggType() string   { return "Invitation" }
func (c *InviteDeclined) EventType() string { return "InviteDeclined" }
