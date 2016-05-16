package redis

import (
	. "github.com/eynstudio/gobreak"
	"github.com/garyburd/redigo/redis"
)

func (p *Redis) Sadd(k string, m ...T) (int, error) {
	return redis.Int(p.Do("SADD", Args(k, m)...))
}
func Sadd(k string, m ...T) (int, error) { return Default.Sadd(k, m...) }

//返回集合存储的key的基数 (集合元素的数量)
func (p *Redis) Scard(k string) (int, error) { return redis.Int(p.Do("SCARD")) }
func Scard(k string) (int, error)            { return Default.Scard(k) }

//返回一个集合与给定集合的差集的元素.
func (p *Redis) Sdiff(k ...interface{}) ([]interface{}, error) {
	return redis.Values(p.Do("SDIFF", k...))
}
func Sdiff(k ...interface{}) ([]interface{}, error) { return Default.Sdiff(k...) }

func (p *Redis) SdiffStore(dest string, k ...interface{}) (int, error) {
	return redis.Int(p.Do("SDIFFSTORE", Args(dest, k)))
}
func SdiffStore(dest string, k ...interface{}) (int, error) { return Default.SdiffStore(dest, k...) }

//返回指定所有的集合的成员的交集.
func (p *Redis) Sinter(k ...interface{}) ([]interface{}, error) {
	return redis.Values(p.Do("SINTER", k...))
}
func Sinter(k ...interface{}) ([]interface{}, error) { return Default.Sinter(k...) }

func (p *Redis) SinterStore(dest string, k ...interface{}) (int, error) {
	return redis.Int(p.Do("SINTERSTORE", Args(dest, k)))
}
func SinterStore(dest string, k ...interface{}) (int, error) { return Default.SinterStore(dest, k...) }

func (p *Redis) Sismember(k string, m T) (bool, error) { return redis.Bool(p.Do("SISMEMBER", k, m)) }
func Sismember(k string, m T) (bool, error)            { return Default.Sismember(k, m) }

func (p *Redis) Smembers(k string) ([]interface{}, error) { return redis.Values(p.Do("SMEMBERS", k)) }
func Smembers(k string, m T) ([]interface{}, error)       { return Default.Smembers(k) }

func (p *Redis) Smove(s, d string, m T) (bool, error) { return redis.Bool(p.Do("SMOVE", s, d, m)) }
func Smove(s, d string, m T) (bool, error)            { return Default.Smove(s, d, m) }

func (p *Redis) Srem(k string, m ...T) (int, error) {
	return redis.Int(p.Do("SREM", Args(k, m)...))
}
func Srem(k string, m ...T) (int, error) { return Default.Srem(k, m...) }

func (p *Redis) Sunion(k ...string) ([]interface{}, error) {
	return redis.Values(p.Do("SUNION", Args(k)...))
}
func Sunion(k ...string) ([]interface{}, error) { return Default.Sunion(k...) }

func (p *Redis) SunionStore(dest string, k ...string) (int, error) {
	return redis.Int(p.Do("SUNIONSTORE", Args(dest, k)...))
}
func SunionStore(dest string, k ...string) (int, error) { return Default.SunionStore(dest, k...) }
