package redis

import (
	. "github.com/eynstudio/gobreak"
	"github.com/garyburd/redigo/redis"
)

type Hash struct {
	Name string
}

func (p *Hash) Exists(id T) (bool, error) {
	i, err := redis.Int(_redis.do("HEXISTS", p.Name, id))
	return i == 1, err
}

func (p *Hash) Get(id T) (interface{}, error) { return _redis.do("HGET", p.Name, id) }
func (p *Hash) GetBytes(id T) ([]byte, error) { return redis.Bytes(_redis.do("HGET", p.Name, id)) }
func (p *Hash) Vals() (interface{}, error)    { return _redis.do("HVALS", p.Name) }

func (p *Hash) Set(id T, t T) (err error) {
	_, err = _redis.do("HSET", p.Name, id, t)
	return
}

func (p *Hash) Del(id T) (err error) {
	_, err = _redis.do("HDEL", p.Name, id)
	return
}

func (p *Redis) Hget(k string, f T) Cmd    { return p.Do("HGET", k, f) }
func (p *Redis) Hset(k string, f, v T) Cmd { return p.Do("HSET", k, f, v) }
