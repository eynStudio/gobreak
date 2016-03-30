package gobreak

import (
	"log"
	"testing"
)

func Test_M(t *testing.T) {
	m := &M{"a": "a", "b": 3}
	log.Println(m.Get("a"))
	log.Println(m.Get("c"))
	log.Println(m)
}
