package redis

import (
	"encoding/json"
	. "github.com/eynstudio/gobreak"
	"github.com/garyburd/redigo/redis"
)

func (p *Redis) Hdel(k string, f ...T) (int, error) {
	return redis.Int(p.Do("HDEL", Args(k, f)...))
}
func Hdel(k string, f ...T) (int, error) { return Default.Hdel(k, f...) }

func (p *Redis) Hexists(k string, f T) (int, error) { return redis.Int(p.Do("HEXISTS", k, f)) }
func Hexists(k string, f T) (int, error)            { return Default.Hexists(k, f) }

func (p *Redis) Hget(k string, f T) (interface{}, error) { return p.Do("HGET", k, f) }
func Hget(k string, f T) (interface{}, error)            { return Default.Hget(k, f) }

func (p *Redis) HgetJson(k string, f T, m T) error {
	data, err := redis.Bytes(p.Hget(k, f))
	NoErrExec(err, func() { err = json.Unmarshal(data, m) })
	return err
}
func HgetJson(k string, f T, m T) error { return Default.HgetJson(k, f, m) }

func (p *Redis) Hgetall(k string) ([]interface{}, error) { return redis.Values(p.Do("HGETALL", k)) }
func Hgetall(k string) ([]interface{}, error)            { return Default.Hgetall(k) }

func (p *Redis) Hincrby(k string, f, inc T) (int, error) { return redis.Int(p.Do("HINCRBY", k, f, inc)) }
func Hincrby(k string, f, inc T) (int, error)            { return Default.Hincrby(k, f, inc) }

func (p *Redis) Hincrbyfloat(k, f string, inc float64) (float64, error) {
	return redis.Float64(p.Do("HINCRBYFLOAT", k, f, inc))
}
func Hincrbyfloat(k, f string, inc float64) (float64, error) { return Default.Hincrbyfloat(k, f, inc) }

func (p *Redis) Hkeys(k string) ([]string, error) { return redis.Strings(p.Do("HKEYS", k)) }
func Hkeys(k string) ([]string, error)            { return Default.Hkeys(k) }

func (p *Redis) Hlen(k string) (int, error) { return redis.Int(p.Do("HLEN", k)) }
func Hlen(k string) (int, error)            { return Default.Hlen(k) }

func (p *Redis) Hmget(k string, args ...T) ([]interface{}, error) {
	return redis.Values(p.Do("HMGET", Args(k, args)...))
}
func Hmget(k string, args ...T) ([]interface{}, error) { return Default.Hmget(k, args...) }

func (p *Redis) Hmset(k string, args ...T) (string, error) {
	return redis.String(p.Do("HMSET", Args(k, args)...))
}
func Hmset(k string, args ...T) (string, error) { return Default.Hmset(k, args...) }

func (p *Redis) Hset(k string, f, v T) (int, error) { return redis.Int(p.Do("HSET", k, f, v)) }
func Hset(k string, f, v T) (int, error)            { return Default.Hset(k, f, v) }

func (p *Redis) HsetJson(k string, f, v T) (int, error) {
	if data, err := json.Marshal(v); err == nil {
		return p.Hset(k, f, data)
	} else {
		return 0, err
	}
}
func HsetJson(k string, f, v T) (int, error) { return Default.HsetJson(k, f, v) }

func (p *Redis) Hvals(k string) ([]interface{}, error) { return redis.Values(p.Do("HVALS", k)) }
func Hvals(k string) ([]interface{}, error)            { return Default.Hvals(k) }
