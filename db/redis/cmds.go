package redis

import (
	. "github.com/eynstudio/gobreak"
	"github.com/garyburd/redigo/redis"
)

func (p *Redis) Hdel(k string, f ...T) (int, error) {
	return redis.Int(p.do("HDEL", buildArgs(k, f)...))
}
func Hdel(k string, f ...T) (int, error)                 { return Default.Hdel(k, f...) }
func (p *Redis) Hexists(k string, f T) (int, error)      { return redis.Int(p.do("HEXISTS", k, f)) }
func Hexists(k string, f T) (int, error)                 { return Default.Hexists(k, f) }
func (p *Redis) Hget(k string, f T) (interface{}, error) { return p.do("HGET", k, f) }
func Hget(k string, f T) (interface{}, error)            { return Default.Hget(k, f) }
func (p *Redis) Hgetall(k string) ([]interface{}, error) { return redis.Values(p.do("HGETALL", k)) }
func Hgetall(k string) ([]interface{}, error)            { return Default.Hgetall(k) }
func (p *Redis) Hincrby(k string, f, inc T) (int, error) { return redis.Int(p.do("HINCRBY", k, f, inc)) }
func Hincrby(k string, f, inc T) (int, error)            { return Default.Hincrby(k, f, inc) }

func (p *Redis) Hincrbyfloat(k, f string, inc float64) (float64, error) {
	return redis.Float64(p.do("HINCRBYFLOAT", k, f, inc))
}
func Hincrbyfloat(k, f string, inc float64) (float64, error) { return Default.Hincrbyfloat(k, f, inc) }

func (p *Redis) Hkeys(k string) ([]string, error) { return redis.Strings(p.do("HKEYS", k)) }
func Hkeys(k string) ([]string, error)            { return Default.Hkeys(k) }
func (p *Redis) Hlen(k string) (int, error)       { return redis.Int(p.do("HLEN", k)) }
func Hlen(k string) (int, error)                  { return Default.Hlen(k) }

func (p *Redis) Hset(k string, f, v T) *Cmd     { return p.Do("HSET", k, f, v) }
func Hset(k string, f, v T) *Cmd                { return Default.Do("HSET", k, f, v) }
func (p *Redis) Hmget(k string, args ...T) *Cmd { return p.Do("HMGET", buildArgs(k, args)...) }
func Hmget(k string, args ...T) *Cmd            { return Default.Do("HMGET", buildArgs(k, args)...) }
func (p *Redis) Hmset(k string, args ...T) *Cmd { return p.Do("HMSET", buildArgs(k, args)...) }
func Hmset(k string, args ...T) *Cmd            { return Default.Do("HMSET", buildArgs(k, args)...) }
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
