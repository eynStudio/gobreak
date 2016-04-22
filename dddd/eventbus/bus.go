package eventbus

import (
	. "github.com/eynstudio/gobreak"
	"log"
)

type Event interface {
	ID() GUID
}

func Publish(event Event) {
	log.Println("publish event ", event)
}
