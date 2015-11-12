//go:generate tpl2str
package main

import (
	"fmt"
)

func main() {
	fmt.Println("gobreak.tools.got!")
	var g Generator
	g.Gen()
}

func loadCfg() {

}
