package redis

import (
	. "github.com/eynstudio/gobreak"
	"github.com/garyburd/redigo/redis"
)

func (p *Redis) Hget(k string, f T) *Cmd        { return p.Do("HGET", k, f) }
func Hget(k string, f T) *Cmd                   { return Default.Do("HGET", k, f) }
func (p *Redis) Hset(k string, f, v T) *Cmd     { return p.Do("HSET", k, f, v) }
func Hset(k string, f, v T) *Cmd                { return Default.Do("HSET", k, f, v) }
func (p *Redis) Hmget(k string, args ...T) *Cmd { return p.Do("HMGET", buildArgs(k, args)...) }
func Hmget(k string, args ...T) *Cmd            { return Default.Do("HMGET", buildArgs(k, args)...) }
func (p *Redis) Hmset(k string, args ...T) *Cmd { return p.Do("HMSET", buildArgs(k, args)...) }
func Hmset(k string, args ...T) *Cmd            { return Default.Do("HMSET", buildArgs(k, args)...) }
func (p *Redis) Hgetall(k string) *Cmd          { return p.Do("HGETALL", k) }
func Hgetall(k string) *Cmd                     { return Default.Do("HGETALL", k) }
func (p *Redis) Hvals(k string) *Cmd            { return p.Do("HVALS", k) }
func Hvals(k string) *Cmd                       { return Default.Do("HVALS", k) }

func (p *Redis) Sadd(k string, args ...T) *Cmd { return p.Do("SADD", buildArgs(k, args)...) }
func Sadd(k string, args ...T) *Cmd            { return Default.Do("SADD", buildArgs(k, args)...) }

func (p *Redis) Zadd(k string, args ...T) *Cmd       { return p.Do("ZADD", buildArgs(k, args)...) }
func Zadd(k string, args ...T) *Cmd                  { return Default.Do("ZADD", buildArgs(k, args)...) }
func (p *Redis) Zincrby(k string, i, m T) *Cmd       { return p.Do("ZINCRBY", k, i, m) }
func Zincrby(k string, i, m T) *Cmd                  { return Default.Do("ZINCRBY", k, i, m) }
func (p *Redis) Zrange(k string, f, t int) *Cmd      { return p.Do("ZRANGE", k, f, t) }
func Zrange(k string, f, t int) *Cmd                 { return Default.Do("ZRANGE", k, f, t) }
func (p *Redis) ZrangeScore(k string, f, t int) *Cmd { return p.Do("ZRANGE", k, f, t, "WITHSCORES") }
func ZrangeScore(k string, f, t int) *Cmd            { return Default.Do("ZRANGE", k, f, t, "WITHSCORES") }

func buildArgs(k T, args ...T) (a redis.Args) {
	a = a.Add(k)
	for _, it := range args {
		a = a.AddFlat(it)
	}
	return
}
