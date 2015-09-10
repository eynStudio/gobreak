package db

type DbCtx interface{
	Shutdown() error
}

