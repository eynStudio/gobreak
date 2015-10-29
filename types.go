package gobreak

type T interface{}

type GUID string

type Entity interface{
	ID() GUID
}

type KeyValue struct {
	K string `K`
	V string `V`
}
type KeyValues struct {
	K string   `K`
	V []string `V`
}
