package redis

import (
	"time"

	"github.com/eynstudio/gobreak/log/datelog"
	"github.com/garyburd/redigo/redis"
)

var Default Redis

func Init(host, pwd string) { Default.Init(host, pwd) }
func Do(cmd string, args ...interface{}) (reply interface{}, err error) {
	return Default.do(cmd, args...)
}

var Log = Default.log.Log
var Logf = Default.log.Logf
var Logln = Default.log.Logln

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
	log  *datelog.DateLog
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

func (p *Redis) Do(cmd string, args ...interface{}) (c *Cmd) {
	rc := p.pool.Get()
	defer rc.Close()
	c = new(Cmd)
	c.reply, c.Err = rc.Do(cmd, args...)
	return
}
