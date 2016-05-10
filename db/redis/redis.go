package redis

import (
	"time"

	"github.com/eynstudio/gobreak/log/datelog"
	"github.com/garyburd/redigo/redis"
)

var _redis Redis

func Init(host, pwd string) { _redis.Init(host, pwd) }
func Do(cmd string, args ...interface{}) (reply interface{}, err error) {
	return _redis.do(cmd, args...)
}

var Log = _redis.log.Log
var Logf = _redis.log.Logf
var Logln = _redis.log.Logln

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
	dbpool *redis.Pool
	host   string
	pwd    string
	log    *datelog.DateLog
}

func (p *Redis) Init(host, pwd string) {
	p.log = datelog.New("./logs/redis")
	p.host = host
	p.pwd = pwd
	p.dbpool = &redis.Pool{
		MaxIdle:     20,
		MaxActive:   10000,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			if c, err := redis.Dial("tcp", host); err != nil {
				return nil, err
			} else {
				c.Do("AUTH", pwd)
				return c, nil
			}
		},
	}
}

func (p *Redis) do(cmd string, args ...interface{}) (reply interface{}, err error) {
	rc := p.dbpool.Get()
	defer rc.Close()
	return rc.Do(cmd, args...)
}
