package gobreak

type T interface{}

type GUID string

func (p GUID) ID() GUID { return p }

type Entity interface {
	ID() GUID
}

type KeyValue struct {
	K string `K`
	V string `V`
}

type KeyValues map[string][]string

func AsKeyValues(lst []KeyValue) (kvs KeyValues) {
	kvs = KeyValues{}
	for _, it := range lst {
		if _, ok := kvs[it.K]; ok {
			kvs[it.K] = append(kvs[it.K], it.V)
		} else {
			kvs[it.K] = []string{it.V}
		}
	}
	return
}
