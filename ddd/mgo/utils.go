package mgo

import (
	. "github.com/eynstudio/gobreak"
	"gopkg.in/mgo.v2/bson"
)

func NewGuid() GUID {
	id := bson.NewObjectId().Hex()
	return GUID(id)
}
