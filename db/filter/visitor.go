package filter

import (
	"fmt"
	"strings"

	. "github.com/eynstudio/gobreak/db"
)

type VisitorResults []SqlArgs

func (p VisitorResults) Join(con string) (vr SqlArgs) {
	if con == "" {
		con = "and"
	}

	var sqls []string
	for _, it := range p {
		if it.Sql != "" {
			sqls = append(sqls, it.Sql)
			vr.Args = append(vr.Args, it.Args...)
		}
	}
	if len(sqls) > 0 {
		vr.Sql = "(" + strings.Join(sqls, " "+con+" ") + ")"
	}
	return
}

type Visitor struct {
}

func (p *Visitor) Visitor(filter Group) SqlArgs {
	return p.VisitGroup(filter)
}

func (p *Visitor) VisitGroup(f Group) (vr SqlArgs) {
	var lst VisitorResults

	ruleResult := p.VisitRules(f.Con, f.Rules)
	lst = append(lst, ruleResult)

	for _, it := range f.Groups {
		lst = append(lst, p.VisitGroup(it))
	}

	return lst.Join(f.Con)
}

func (p *Visitor) VisitRules(con string, rules []Rule) (vr SqlArgs) {
	var lst VisitorResults
	for _, it := range rules {
		lst = append(lst, p.VisitRule(it))
	}
	return lst.Join(con)
}

func (p *Visitor) VisitRule(f Rule) (vr SqlArgs) {
	if f.O == "like" {
		vr.Sql = fmt.Sprintf("%s %s ?", f.F, f.O)
		vr.Args = append(vr.Args, "%"+f.V1+"%")
	} else if f.O == "=" {
		vr.Sql = fmt.Sprintf("%s %s ?", f.F, f.O)
		vr.Args = append(vr.Args, f.V1)
	} else if f.O == "in" {
		var ss []string
		vlst := strings.Split(f.V1, ",")
		for _, it := range vlst {
			ss = append(ss, "?")
			vr.Args = append(vr.Args, it)
		}
		vs := strings.Join(ss, ",")
		vr.Sql = fmt.Sprintf("%s %s (%s)", f.F, f.O, vs)
	}
	return
}