//go:generate tpl2str
package main

import (
	"fmt"
)

func main() {
	fmt.Println("start gobreak.tools.got")
	var g Generator
	g.loadCfg()
	g.Gen()
}
