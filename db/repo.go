package db

import (
	. "github.com/eynstudio/gobreak"
)

type Repo interface {
	All(m T) T
	Get(id T, m T) T
	Save(id T, m T)
	Del(id T)
}

type BaseRepo interface {
	Repo
}
