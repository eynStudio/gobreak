package db

type DbCtx interface {
	GetCollection(name string) interface{}
	Shutdown() error
}
