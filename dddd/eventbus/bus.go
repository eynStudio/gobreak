package eventbus

import (
	"log"

	. "github.com/eynstudio/gobreak"
)

func Publish(event Event) {
	log.Println("publish event ", event)
}
