package redis

import (
	"bytes"
	"encoding/json"
	"reflect"

	. "github.com/eynstudio/gobreak"
	"github.com/garyburd/redigo/redis"
)

type JsonHash struct {
	Hash
}

func (p *JsonHash) Get(id string, t T) (has bool, err error) {
	data, err := p.GetBytes(id)
	switch err {
	case nil:
		return true, json.Unmarshal(data, &t)
	case redis.ErrNil:
		return false, nil
	default:
		return false, err
	}
}

func (p *JsonHash) Vals(t T) error {
	if data, err := redis.ByteSlices(_redis.do("HVALS", p.Name)); err != nil {
		return err
	} else {
		resultv := reflect.ValueOf(t)
		slicev := resultv.Elem()
		elem := reflect.TypeOf(t).Elem().Elem()
		for _, it := range data {
			itm := reflect.New(elem)
			json.Unmarshal(it, itm.Interface())
			slicev = reflect.Append(slicev, itm.Elem())
		}
		resultv.Elem().Set(slicev.Slice(0, slicev.Len()))
		return nil
	}
}

func (p *JsonHash) Set(id string, t T) error {
	if data, err := json.Marshal(t); err != nil {
		return err
	} else {
		return p.Set(id, data)
	}
}

func (p *JsonHash) SetIfDiff(id string, t T) (changed bool, e error) {
	if ok, err := p.Exists(id); err != nil {
		return false, err
	} else if !ok {
		return true, p.Set(id, t)
	}

	var data, data2 []byte
	var err error
	data, err = p.GetBytes(id)
	NoErrExec(err, func() { data2, err = json.Marshal(t) })
	NoErrExec(err, func() {
		if !bytes.Equal(data, data2) {
			err = p.Hash.Set(id, data2)
			changed = true
		}
	})
	return changed, err
}
