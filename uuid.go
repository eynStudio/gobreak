package gobreak

import (
	"github.com/satori/go.uuid"
)

func Uuid0() uuid.UUID                          { return uuid.Nil }
func Uuid1() uuid.UUID                          { return uuid.NewV1() }
func Uuid3(ns uuid.UUID, name string) uuid.UUID { return uuid.NewV3(ns, name) }
func Uuid4() uuid.UUID                          { return uuid.NewV4() }
func Uuid5(ns uuid.UUID, name string) uuid.UUID { return uuid.NewV5(ns, name) }
