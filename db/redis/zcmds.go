package redis

import (
	. "github.com/eynstudio/gobreak"
	"github.com/garyburd/redigo/redis"
)

func (p *Redis) Zadd(k string, args ...T) (int, error) {
	return redis.Int(p.Do("ZADD", Args(k, args)...))
}
func Zadd(k string, args ...T) (int, error) { return Default.Zadd(k, args...) }

//返回key的有序集元素个数。
func (p *Redis) Zcard(k string) (int, error) { return redis.Int(p.Do("ZCARD")) }
func Zcard(k string) (int, error)            { return Default.Zcard(k) }

//返回有序集key中，score值在min和max之间(默认包括score值等于min或max)的成员。
func (p *Redis) Zcount(k string, min, max T) (int, error) {
	return redis.Int(p.Do("ZCOUNT", k, min, max))
}
func Zcount(k string, min, max T) (int, error) { return Default.Zcount(k, min, max) }

func (p *Redis) Zincrby(k string, inc, m T) (interface{}, error) { return p.Do("ZINCRBY", k, inc, m) }
func Zincrby(k string, inc, m T) (interface{}, error)            { return Default.Zincrby(k, inc, m) }

func (p *Redis) ZincrbyInt(k string, inc int, m T) (int, error) {
	return redis.Int(p.Zincrby(k, inc, m))
}
func ZincrbyInt(k string, inc int, m T) (int, error) { return Default.ZincrbyInt(k, inc, m) }

func (p *Redis) ZincrbyF(k string, inc, m T) (float64, error) {
	return redis.Float64(p.Zincrby(k, inc, m))
}
func ZincrbyF(k string, inc, m T) (float64, error) { return Default.ZincrbyF(k, inc, m) }

func (p *Redis) Zinterstore(d string, n int, keys []string, agg string, w []float64) (int, error) {
	return redis.Int(p.Do("ZINTERSTORE", Args(d, n, keys, agg, w)...))
}
func Zinterstore(d string, n int, keys []string, agg string, w []float64) (int, error) {
	return Default.Zadd(d, n, keys, agg, w)
}

func (p *Redis) Zrem(k string, m ...T) (int, error) { return redis.Int(p.Do("ZREM", Args(k, m)...)) }
func Zrem(k string, m ...T) (int, error)            { return Default.Zrem(k, m...) }
