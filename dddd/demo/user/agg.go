package user

import (
	. "github.com/eynstudio/gobreak"
	. "github.com/eynstudio/gobreak/dddd"
	"github.com/eynstudio/gobreak/dddd/cmdbus"
	"log"
)

func init() {
	log.Println("init user...")

}

type User struct {
	Name string
	Age  int
}

type UserAgg struct {
	AggBase
}

type SaveUser struct{}
type UpdateAge struct{}
type UserSaved struct{}
type UserUpdated struct{}

func (p *UserAgg) Root() *User { return p.AggBase.Root.(*User) }

func (p User) ID() GUID                       { return "" }
func (p SaveUser) ID() GUID                   { return "" }
func (p *UpdateAge) ID() GUID                 { return "" }
func (p *UserSaved) ID() GUID                 { return "" }
func (p *UserUpdated) ID() GUID               { return "" }
func (p *UserAgg) RegistedCmds() []cmdbus.Cmd { return []cmdbus.Cmd{&SaveUser{}, &UpdateAge{}} }

func (p *UserAgg) HandleCmd(cmd cmdbus.Cmd) error {
	log.Println(cmd)
	switch cmd := cmd.(type) {
	case *SaveUser:
		p.ApplyEvent((*UserSaved)(cmd))
	case *UpdateAge:
		p.ApplyEvent((*UserUpdated)(cmd))
	default:
		log.Println("UserAgg HandleCmd: no handler")
	}
	return nil
}

func (p *UserAgg) ApplyEvent(events Event) {
	log.Println("apply event ", events)
	p.StoreEvent(events)
}
