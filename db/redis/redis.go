package redis

import (
	. "github.com/eynstudio/gobreak"
	"time"

	"github.com/garyburd/redigo/redis"
)

var _redis Redis

func Init(server, pwd string) { _redis.Init(server, pwd) }
func Do(cmd string, args ...interface{}) (reply interface{}, err error) {
	return _redis.do(cmd, args...)
}

var Int = redis.Int
var Int64 = redis.Int64
var Uint64 = redis.Uint64

var Float64 = redis.Float64
var String = redis.String
var Bytes = redis.Bytes
var Bool = redis.Bool
var Values = redis.Values
var Strings = redis.Strings
var ByteSlices = redis.ByteSlices
var Ints = redis.Ints
var StringMap = redis.StringMap
var IntMap = redis.IntMap
var Int64Map = redis.Int64Map

type Redis struct {
	pool *redis.Pool
}

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func (p *Redis) Init(server, pwd string) {
	p.pool = newPool(server, pwd)
}

func (p *Redis) do(cmd string, args ...interface{}) (reply interface{}, err error) {
	rc := p.pool.Get()
	defer rc.Close()
	return rc.Do(cmd, args...)
}

func (p *Redis) Do(cmd string, args ...interface{}) (c Cmd) {
	rc := p.pool.Get()
	defer rc.Close()
	c.reply, c.Err = rc.Do(cmd, args...)
	return
}

type Cmd struct {
	Error
	reply T
}

func (p *Cmd) Bytes() (bytes []byte) {
	bytes, p.Err = redis.Bytes(p.reply, p.Err)
	return
}

func (p *Cmd) Int() (i int) {
	i, p.Err = redis.Int(p.reply, p.Err)
	return
}
