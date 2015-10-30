package gobreak

type T interface{}

type GUID string
func (p GUID) ID() GUID {return p}

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
