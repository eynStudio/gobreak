package main

import (
	"github.com/eynstudio/gobreak/dddd/cmdbus"
	"github.com/eynstudio/gobreak/dddd/demo/user"

	"log"
)

func main() {
	log.Println("dddd demo")

	cmdbus.Exec(user.SaveUser{})

	log.Println("dddd demo end")
}
