package di2

import "github.com/eynstudio/gobreak/log"

var g *Graph = New()

func New() *Graph                              { return &Graph{} }
func Reg(values ...interface{}) error          { return g.Reg(values...) }
func RegAs(name string, val interface{}) error { return g.RegAs(name, val) }
func SetLog(l log.ILogger)                     { g.Logger = l }
