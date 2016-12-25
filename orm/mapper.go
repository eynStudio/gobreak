package orm

import "github.com/eynstudio/gox/x/strx"

type MapperFn func(string) string

type Mapper interface {
	Map2Db(string) string
	Map2Obj(string) string
}

type baseMapper struct {
	map2Db  MapperFn
	map2Obj MapperFn
}

func New(map2db, map2obj MapperFn) Mapper {
	return &baseMapper{map2Db: map2db, map2Obj: map2obj}
}
func (mp *baseMapper) Map2Db(s string) string  { return mp.map2Db(s) }
func (mp *baseMapper) Map2Obj(s string) string { return mp.map2Obj(s) }

func MapperSelf() Mapper    { return New(map2Self, map2Self) }
func Mapper_() Mapper       { return New(strx.UnderScoreCase, strx.UpperCamel) }
func MapperUpCamel() Mapper { return New(strx.UpperCamel, strx.UpperCamel) }

func map2Self(n string) string { return n }
