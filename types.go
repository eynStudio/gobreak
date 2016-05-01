package gobreak

type T interface{}

type GUID string

func (p GUID) ID() GUID       { return p }
func (p GUID) String() string { return string(p) }
func (p GUID) IsEmpty() bool  { return len(p) == 0 }

type Entity interface {
	ID() GUID
}

type KeyValue struct {
	K string `K`
	V string `V`
}

type Params []KeyValue

func (p Params) Has(k, v string) bool {
	for _, it := range p {
		if it.K == k && it.V == v {
			return true
		}
	}
	return false
}

func (p Params) HasKey(key string) bool {
	for _, it := range p {
		if it.K == key {
			return true
		}
	}
	return false
}

func (p Params) Get(key string) *KeyValue {
	for _, it := range p {
		if it.K == key {
			return &it
		}
	}
	return nil
}

func (p Params) GetAll(key string) (lst []KeyValue) {
	for _, it := range p {
		if it.K == key {
			lst = append(lst, it)
		}
	}
	return
}

func (p Params) GetValue(key string) string {
	for _, it := range p {
		if it.K == key {
			return it.V
		}
	}
	return ""
}

func (p Params) GetValues(key string) (strs []string) {
	for _, it := range p {
		if it.K == key {
			strs = append(strs, it.V)
		}
	}
	return
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
