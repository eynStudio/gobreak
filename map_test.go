package gobreak

import (
	"log"
	"testing"
)

func Test_M(t *testing.T) {
	m := &M{"a": "a", "b": 3}
	log.Println(m.GetStr("a"))
	log.Println(m.GetInt("c"))
	log.Println(m)
}
