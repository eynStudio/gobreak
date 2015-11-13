package {{.ParentPkg}}

import (
	"fmt"
	. "{{.ParentPath}}/{{.Pkg}}"

	. "github.com/eynstudio/gobreak"
	. "github.com/eynstudio/gobreak/db/mgo"
	. "github.com/eynstudio/gobreak/ddd"
	"github.com/eynstudio/gobreak/orm"
	"gopkg.in/mgo.v2/bson"
)

type {{.Name}}Repo struct {
	MgoRepo
}

func New{{.Name}}Repo() *{{.Name}}Repo {
	return &{{.Name}}Repo{NewMgoRepo("{{.Repo}}", func() T { return &{{.Name}}{} })}
}

type {{.Name}}EventHandler struct {
	Repo *{{.Name}}Repo  `di`
	Orm  *orm.Orm `di`
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
