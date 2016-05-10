package redis

import (
	. "github.com/eynstudio/gobreak"
	"github.com/garyburd/redigo/redis"
)

type Hash struct {
	Name string
}

func (p *Hash) Exists(id T) (bool, error) {
	i, err := redis.Int(Default.do("HEXISTS", p.Name, id))
	return i == 1, err
}

func (p *Hash) Get(id T) (interface{}, error) { return Default.do("HGET", p.Name, id) }
func (p *Hash) GetBytes(id T) ([]byte, error) { return redis.Bytes(Default.do("HGET", p.Name, id)) }
func (p *Hash) Vals() (interface{}, error)    { return Default.do("HVALS", p.Name) }

func (p *Hash) Set(id T, t T) (err error) {
	_, err = Default.do("HSET", p.Name, id, t)
	return
}

func (p *Hash) Del(id T) (err error) {
	_, err = Default.do("HDEL", p.Name, id)
	return
}

func (p *Redis) Hget(k string, f T) *Cmd        { return p.Do("HGET", k, f) }
func (p *Redis) Hset(k string, f, v T) *Cmd     { return p.Do("HSET", k, f, v) }
func (p *Redis) Hmget(k string, args ...T) *Cmd { return p.Do("HMGET", buildArgs(k, args)...) }
func (p *Redis) Hmset(k string, args ...T) *Cmd { return p.Do("HMSET", buildArgs(k, args)...) }
func (p *Redis) Hgetall(k string) *Cmd          { return p.Do("HGETALL", k) }
func (p *Redis) Hvals(k string) *Cmd            { return p.Do("HVALS", k) }

func buildArgs(k T, args ...T) (a redis.Args) {
	a = a.Add(k)
	for _, it := range args {
		a = a.AddFlat(it)
	}
	return
}
