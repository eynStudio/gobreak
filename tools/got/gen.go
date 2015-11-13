package main

import (
	"bytes"
	//	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	//	"os"
	"text/template"

	"github.com/eynstudio/gobreak"
)

type Generator struct {
	buf bytes.Buffer
	cfg *CfgVm
}

func (p *Generator) loadCfg() {
	var cfg Cfg
	if !gobreak.LoadJson("got.json", &cfg) {
		panic("Can not load got.json file")
	}
	p.cfg = cfg.parseCfgVm()
}

func (p *Generator) Gen() {
	for _, c := range p.cfg.Aggs {
		p.GenAgg(&c)
	}
}

func (p *Generator) GenAgg(agg *AggVm) {
	os.Mkdir(strings.ToLower(agg.Name), os.ModeDir)
	p.GenFile(agg, tpl_agg, agg.Name+"/agg.go")
	p.GenFile(agg, tpl_cmds, agg.Name+"/cmds.go")
	p.GenFile(agg, tpl_events, agg.Name+"/events.go")
	p.GenFile(agg, tpl_entities, agg.Name+"/entities.go")
	p.GenFile(agg, tpl_read, agg.Pkg+"_read.go")
}

func (p *Generator) GenFile(agg *AggVm, tpl, tofile string) {
	var buf bytes.Buffer
	tmpl, err := template.New("agg").Parse(tpl)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(&buf, agg)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(tofile, buf.Bytes(), 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}
