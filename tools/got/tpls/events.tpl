package {{.Pkg}}

import (
	. "github.com/eynstudio/gobreak"
	. "github.com/eynstudio/gobreak/ddd"
)

type {{.Name}}Saved {{.Name}}

func (p *{{.Name}}Saved) ID() GUID { return p.Id }

type {{.Name}}Deleted IdEvent

func (p *{{.Name}}Deleted) ID() GUID { return p.Id }