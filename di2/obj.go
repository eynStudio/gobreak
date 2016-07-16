package di2

import (
	"bytes"
	"fmt"
	"reflect"
)

type Obj struct {
	val      interface{}
	name     string //for named
	fields   map[string]*Obj
	diType   reflect.Type
	diVal    reflect.Value
	private  bool
	created  bool
	embedded bool
}

func (p *Obj) setField(field string, dep *Obj) {
	if p.fields == nil {
		p.fields = make(map[string]*Obj)
	}
	p.fields[field] = dep
}

func (p *Obj) string() string {
	var buf bytes.Buffer
	fmt.Fprint(&buf, p.diType)
	if p.name != "" {
		fmt.Fprintf(&buf, " named %s", p.name)
	}
	return buf.String()
}
