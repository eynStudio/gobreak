package user

import (
	"log"

	. "github.com/eynstudio/gobreak"
	. "github.com/eynstudio/gobreak/dddd/ddd"
)

type User struct {
	Id   GUID   `bson:"_id,omitempty"`
	Name string `Name`
	Age  int    `Age`
}

type UserAgg struct {
	AggBase
	root User
}

type SaveUser User
type UpdateAge struct {
	Id  GUID
	Age int
}
type UserSaved User
type UserUpdated UpdateAge

func (p *UserAgg) Root() Entity { return &p.root }

func (p User) ID() GUID    { return p.Id }
func (p UserAgg) ID() GUID { return p.Root().ID() }

func (p SaveUser) ID() GUID            { return p.Id }
func (p *UpdateAge) ID() GUID          { return p.Id }
func (p *UserSaved) ID() GUID          { return p.Id }
func (p *UserUpdated) ID() GUID        { return p.Id }
func (p *UserAgg) RegistedCmds() []Cmd { return []Cmd{&SaveUser{}, &UpdateAge{}} }

func (p *UserAgg) HandleCmd(cmd Cmd) error {
	log.Println("useragg handle cmd", cmd)
	switch cmd := cmd.(type) {
	case *SaveUser:
		p.root = User(*cmd)
		p.ApplyEvent((*UserSaved)(cmd))
	case *UpdateAge:
		p.ApplyEvent((*UserUpdated)(cmd))
	default:
		log.Println("UserAgg HandleCmd: no handler")
	}
	return nil
}

func (p *UserAgg) ApplyEvent(events Event) {
	log.Println("apply event ", events, p.root)
	p.StoreEvent(events)
}
