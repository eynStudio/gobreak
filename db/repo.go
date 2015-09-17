package db

type Repo interface {
	All(m interface{}) interface{}
	Get(id interface{}, m interface{}) interface{}
	Save(id interface{}, m interface{})
}

type BaseRepo interface {
	Repo
}
