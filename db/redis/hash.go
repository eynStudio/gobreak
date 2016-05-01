package redis

import (
	"bytes"
	"encoding/json"

	. "github.com/eynstudio/gobreak"
	"github.com/garyburd/redigo/redis"
)

type Hash struct {
	Name string
}

func (p *Hash) GetBytes(id string) ([]byte, error) {
	return redis.Bytes(_redis.do("HGET", p.Name, id))
}

func (p *Hash) Exists(id string) bool {
	i, _ := redis.Int(_redis.do("HEXISTS", p.Name, id))
	return i == 1
}

func (p *Hash) Get(id string) (interface{}, error) { return _redis.do("HGET", p.Name, id) }
func (p *Hash) Set(id string, t T) (interface{}, error) {
	return _redis.do("HSET", p.Name, id, t)
}
func (p *Hash) Del(id string) (interface{}, error) { return _redis.do("HDEL", p.Name, id) }

type JsonHash struct {
	Hash
}

func (p *JsonHash) Get(id string, t T) error {
	if data, err := p.GetBytes(id); err != nil {
		return err
	} else {
		return json.Unmarshal(data, &t)
	}
}

func (p *JsonHash) Set(id string, t T) (interface{}, error) {
	if data, err := json.Marshal(t); err == nil {
		return _redis.do("HSET", p.Name, id, data)
	} else {
		return nil, err
	}
}

func (p *JsonHash) SetIfDiff(id string, t T) bool {
	if !p.Exists(id) {
		p.Set(id, t)
		return true
	}

	var data, data2 []byte
	var err error
	if data, err = p.GetBytes(id); err != nil {
		return false
	}

	if data2, err = json.Marshal(t); err != nil {
		if !bytes.Equal(data, data2) {
			_, err = _redis.do("HSET", p.Name, id, data2)
			return true
		}
	}
	return false
}
