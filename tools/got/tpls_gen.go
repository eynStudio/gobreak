package main

const (
tpl_agg=`package {{.Pkg}}

import (
	"fmt"

	. "github.com/eynstudio/gobreak"
	. "github.com/eynstudio/gobreak/ddd"
)


type {{.AggName}} struct {
	*AggregateBase
	StateModel {{.Name}}
}

func New{{.AggName}}(id GUID) Aggregate {
	return &{{.AggName}}{
		AggregateBase: NewAggregateBase(id),
	}
}

func (p *{{.AggName}}) RegistedCmds() []Cmd {
	return []Cmd{&Save{{.Name}}{}, &Del{{.Name}}{}}
}

func (p *{{.AggName}}) HandleCmd(cmd Cmd) error {
	switch cmd := cmd.(type) {
	case *Save{{.Name}}:
		p.ApplyEvent((*{{.Name}}Saved)(cmd))
	case *Del{{.Name}}:
		p.ApplyEvent((*{{.Name}}Deleted)(cmd))
	default:
		fmt.Println("{{.AggName}} HandleCmd: no handler")
	}
	return nil
}

func (p *{{.AggName}}) ApplyEvent(event Event) {
	switch evt := event.(type) {
	case *{{.Name}}Saved:
		p.StateModel = {{.Name}}(*evt)
	case *{{.Name}}Deleted:
		p.StateModel = {{.Name}}{}
	}
	p.IncrementVersion()
	p.StoreEvent(event)
}

func (p *{{.AggName}}) GetSnapshot() T {
	return &p.StateModel
}
`
tpl_cmds=`package {{.Pkg}}

import (
	. "github.com/eynstudio/gobreak"
	. "github.com/eynstudio/gobreak/ddd"
)

type Save{{.Name}} {{.Name}}

func (p *Save{{.Name}}) ID() GUID { return p.Id }

type Del{{.Name}} IdCmd

func (p *Del{{.Name}}) ID() GUID { return p.Id }`
tpl_controller=`package {{.Name}}

import (
	. "github.com/eynstudio/goweb/mgo"
)`
tpl_entities=`package {{.Pkg}}

import (
	. "github.com/eynstudio/gobreak"
	"github.com/eynstudio/gobreak/db/mgo"
)

type {{.Name}} struct {
	Id     GUID   ` + "`" + `bson:"_id,omitempty"` + "`" + `
	UserId GUID   ` + "`" + `UserId` + "`" + `
}

func New{{.Name}}(uid GUID) *{{.Name}}{
	return &{{.Name}}{
		Id:mgo.NewGuid(),
		UserId:uid,
	}
}`
tpl_events=`package {{.Pkg}}

import (
	. "github.com/eynstudio/gobreak"
	. "github.com/eynstudio/gobreak/ddd"
)

type {{.Name}}Saved {{.Name}}

func (p *{{.Name}}Saved) ID() GUID { return p.Id }

type {{.Name}}Deleted IdEvent

func (p *{{.Name}}Deleted) ID() GUID { return p.Id }`
tpl_read=`package {{.ParentPkg}}

import (
	"fmt"
	. "{{.ParentPath}}/{{.Pkg}}"

	. "github.com/eynstudio/gobreak"
	. "github.com/eynstudio/gobreak/db/mgo"
	. "github.com/eynstudio/gobreak/ddd"
	"github.com/eynstudio/gobreak/orm"
)

type {{.Name}}Repo struct {
	MgoRepo
}

func New{{.Name}}Repo() *{{.Name}}Repo {
	return &{{.Name}}Repo{NewMgoRepo("{{.Repo}}", func() T { return &{{.Name}}{} })}
}

type {{.Name}}EventHandler struct {
	Repo *{{.Name}}Repo  ` + "`" + `di` + "`" + `
	Orm  *orm.Orm ` + "`" + `di` + "`" + `
}

func (p *{{.Name}}EventHandler) RegistedEvents() []Event {
	return []Event{&{{.Name}}Saved{}, &{{.Name}}Deleted{}}
}

func (p *{{.Name}}EventHandler) HandleEvent(event Event) {
	switch event := event.(type) {
	case *{{.Name}}Saved:
		p.Repo.Save(event.ID(), event)
	case *{{.Name}}Deleted:
		p.Repo.Del(event.ID())
	default:
		fmt.Println("{{.Name}}EventHandler: no handler")
	}
}
`

)
