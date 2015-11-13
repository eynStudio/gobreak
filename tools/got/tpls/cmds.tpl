package {{.Pkg}}

import (
	. "github.com/eynstudio/gobreak"
	. "github.com/eynstudio/gobreak/ddd"
)

type Save{{.Name}} {{.Name}}

func (p *Save{{.Name}}) ID() GUID { return p.Id }

type Del{{.Name}} IdCmd

func (p *Del{{.Name}}) ID() GUID { return p.Id }