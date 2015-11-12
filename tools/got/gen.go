package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Generator struct {
	buf bytes.Buffer
}

func (p *Generator) Gen() {
	p.Printf("data")
	os.Mkdir("abc", os.ModeDir)
	err := ioutil.WriteFile("abc/abc.go", p.buf.Bytes(), 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

func (p *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&p.buf, format, args...)
}
