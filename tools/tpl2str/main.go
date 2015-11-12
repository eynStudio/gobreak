package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

func main() {
	var g Generator
	g.Gen()
}

type Generator struct {
	buf bytes.Buffer
}

func (p *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&p.buf, format, args...)
}

func (p *Generator) Gen() {
	files, _ := filepath.Glob("tpls/*")
	fmt.Println(files)
	p.Printf("package main\n")
	p.Printf("\nconst (\n")
	for _, f := range files {
		p.genFile(f)
	}
	p.Printf("\n)\n")

	err := ioutil.WriteFile("tpls_gen.go", p.buf.Bytes(), 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

func (p *Generator) genFile(f string) {
	bytes, _ := ioutil.ReadFile(f)
	f = strings.TrimLeft(f, "tpls\\")
	f = strings.TrimRight(f, ".tpl")
	p.Printf("tpl_%s=`%s`", f, bytes)
}
