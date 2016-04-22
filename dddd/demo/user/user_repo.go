package user

import (
	. "github.com/eynstudio/gobreak"
	. "github.com/eynstudio/gobreak/db/mgo"
	"github.com/eynstudio/gobreak/dddd"
	"github.com/eynstudio/gobreak/di"
)

func Init() {
	up := NewUserRepo()
	di.Map(up).Apply(up.MgoRepo)
	dddd.Reg(&UserAgg{}, up)
}

type UserRepo struct {
	MgoRepo
}

func NewUserRepo() *UserRepo {
	return &UserRepo{NewMgoRepo("Eyn", func() T { return &User{} })}
}
