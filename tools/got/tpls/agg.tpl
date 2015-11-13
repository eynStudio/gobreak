package {{.Pkg}}

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
