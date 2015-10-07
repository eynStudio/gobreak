package db

import(
		. "github.com/eynstudio/gobreak"
)

type DbCtx interface {
	GetCollection(name string) T
	Shutdown() error
}
