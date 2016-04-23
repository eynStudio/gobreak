package main

import (
	"github.com/eynstudio/gobreak/dddd/cmdbus"
	"github.com/eynstudio/gobreak/dddd/demo/user"

	. "github.com/eynstudio/gobreak"
	. "github.com/eynstudio/gobreak/db/mgo"
	"github.com/eynstudio/gobreak/di"

	"log"
)

func main() {
	log.Println("dddd demo")

	mgoCtx := NewMgoCtx(&MgoCfg{})
	di.Map(mgoCtx)
	user.Init()

	di.Invoke(run)

	log.Println("dddd demo end")
}

func run(uRepo *user.UserRepo) {
	log.Println("run...", uRepo)
	cmdbus.Exec(&user.SaveUser{Id: "eyn2", Age: 333, Name: "eeee"})

	log.Printf("%#v", uRepo.All())
}
