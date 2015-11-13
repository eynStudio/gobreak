package main

import (
	"strings"
)

type Cfg struct {
	Package string
	Path    string
	Aggs    []Agg
}

type Agg struct {
	Name string
	Repo string
}

type CfgVm struct {
	Pkg  string
	Path string
	Aggs []AggVm
}
type AggVm struct {
	ParentPkg  string
	ParentPath string
	Pkg        string
	Name       string
	AggName    string
	Repo       string
}

func (p *Cfg) parseCfgVm() *CfgVm {
	var cfg CfgVm
	cfg.Pkg = p.Package
	cfg.Path = p.Path

	for _, a := range p.Aggs {
		var vm AggVm
		vm.ParentPkg = p.Package
		vm.ParentPath = p.Path
		vm.Pkg = a.Name
		vm.Name = strings.Title(a.Name)
		vm.AggName = vm.Name + "Agg"
		vm.Repo = a.Repo

		cfg.Aggs = append(cfg.Aggs, vm)
	}
	return &cfg
}
