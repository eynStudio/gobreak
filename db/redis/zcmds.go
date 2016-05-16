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

// agg is min,max, or others means sum
// when len(w) <> len(keys), not set WEIGHTS
func (p *Redis) Zinterstore(d string, keys []string, agg string, w []float64) (int, error) {
	args := Args(d, len(keys), keys)
	if agg == "MIN" || agg == "MAX" {
		args.Add("AGGREGATE", agg)
	}
	if len(w) > 0 && len(keys) == len(w) {
		args.Add("WEIGHTS").AddFlat(w)
	}
	return redis.Int(p.Do("ZINTERSTORE", args...))
}
func Zinterstore(d string, keys []string, agg string, w []float64) (int, error) {
	return Default.Zinterstore(d, keys, agg, w)
}

func (p *Redis) Zrange(k string, f, t int) ([]interface{}, error) {
	return redis.Values(p.Do("ZRANGE", k, f, t))
}
func Zrange(k string, f, t int) ([]interface{}, error) { return Default.Zrange(k, f, t) }

func (p *Redis) ZrangeScore(k string, f, t int) ([]interface{}, error) {
	return redis.Values(p.Do("ZRANGE", k, f, t, "WITHSCORES"))
}
func ZrangeScore(k string, f, t int) ([]interface{}, error) { return Default.ZrangeScore(k, f, t) }

func (p *Redis) ZrangeByLex(k string, min, max T, useLimit bool, offset, count int) ([]interface{}, error) {
	args := Args(k, min, max)
	if useLimit {
		args.Add("LIMIT", offset, count)
	}
	return redis.Values(p.Do("ZRANGEBYLEX", args...))
}
func ZrangeByLex(k string, min, max T, useLimit bool, offset, count int) ([]interface{}, error) {
	return Default.ZrangeByLex(k, min, max, useLimit, offset, count)
}

func (p *Redis) ZrangeByScore(k string, min, max T, useLimit bool, offset, count int) ([]interface{}, error) {
	args := Args(k, min, max)
	if useLimit {
		args.Add("LIMIT", offset, count)
	}
	return redis.Values(p.Do("ZRANGEBYSCORE", args...))
}
func ZrangeByScore(k string, min, max T, useLimit bool, offset, count int) ([]interface{}, error) {
	return Default.ZrangeByScore(k, min, max, useLimit, offset, count)
}

func (p *Redis) ZrangeByScoreWiScore(k string, min, max T, useLimit bool, offset, count int) ([]interface{}, error) {
	args := Args(k, min, max, "WITHSCORES")
	if useLimit {
		args.Add("LIMIT", offset, count)
	}
	return redis.Values(p.Do("ZRANGEBYSCORE", args...))
}
func ZrangeByScoreWiScore(k string, min, max T, useLimit bool, offset, count int) ([]interface{}, error) {
	return Default.ZrangeByScoreWiScore(k, min, max, useLimit, offset, count)
}

//返回有序集key中成员member的排名
func (p *Redis) Zrank(k string, m T) (int, error) { return redis.Int(p.Do("ZRANK", k, m)) }
func Zrank(k string, m T) (int, error)            { return Default.Zrank(k, m) }

func (p *Redis) Zrem(k string, m ...T) (int, error) { return redis.Int(p.Do("ZREM", Args(k, m)...)) }
func Zrem(k string, m ...T) (int, error)            { return Default.Zrem(k, m...) }

func (p *Redis) ZrevRange(k string, f, t int) ([]interface{}, error) {
	return redis.Values(p.Do("ZREVRANGE", k, f, t))
}
func ZrevRange(k string, f, t int) ([]interface{}, error) { return Default.ZrevRange(k, f, t) }

func (p *Redis) ZrevRangeScore(k string, f, t int) ([]interface{}, error) {
	return redis.Values(p.Do("ZREVRANGE", k, f, t, "WITHSCORES"))
}
func ZrevRangeScore(k string, f, t int) ([]interface{}, error) { return Default.ZrevRangeScore(k, f, t) }

func (p *Redis) ZrevRangeByLex(k string, min, max T, useLimit bool, offset, count int) ([]interface{}, error) {
	args := Args(k, min, max)
	if useLimit {
		args.Add("LIMIT", offset, count)
	}
	return redis.Values(p.Do("ZREVRANGEBYLEX", args...))
}
func ZrevRangeByLex(k string, min, max T, useLimit bool, offset, count int) ([]interface{}, error) {
	return Default.ZrevRangeByLex(k, min, max, useLimit, offset, count)
}

func (p *Redis) ZrevRangeByScore(k string, min, max T, useLimit bool, offset, count int) ([]interface{}, error) {
	args := Args(k, min, max)
	if useLimit {
		args.Add("LIMIT", offset, count)
	}
	return redis.Values(p.Do("ZREVRANGEBYSCORE", args...))
}
func ZrevRangeByScore(k string, min, max T, useLimit bool, offset, count int) ([]interface{}, error) {
	return Default.ZrevRangeByScore(k, min, max, useLimit, offset, count)
}

func (p *Redis) ZrevRangeByScoreWiScore(k string, min, max T, useLimit bool, offset, count int) ([]interface{}, error) {
	args := Args(k, min, max, "WITHSCORES")
	if useLimit {
		args.Add("LIMIT", offset, count)
	}
	return redis.Values(p.Do("ZREVRANGEBYSCORE", args...))
}
func ZrevRangeByScoreWiScore(k string, min, max T, useLimit bool, offset, count int) ([]interface{}, error) {
	return Default.ZrevRangeByScoreWiScore(k, min, max, useLimit, offset, count)
}

//返回有序集key中成员member的排名
func (p *Redis) ZrevRank(k string, m T) (int, error) { return redis.Int(p.Do("ZREVRANK", k, m)) }
func ZrevRank(k string, m T) (int, error)            { return Default.ZrevRank(k, m) }

//返回有序集key中，成员member的score值。
func (p *Redis) Zscore(k string, m T) (interface{}, error) { return p.Do("ZSCORE", k, m) }
func Zscore(k string, m T) (interface{}, error)            { return Default.Zscore(k, m) }
func (p *Redis) ZscoreInt(k string, m T) (int, error)      { return redis.Int(p.Do("ZSCORE", k, m)) }
func ZscoreInt(k string, m T) (int, error)                 { return Default.ZscoreInt(k, m) }
func (p *Redis) ZscoreF64(k string, m T) (float64, error)  { return redis.Float64(p.Do("ZSCORE", k, m)) }
func ZscoreF64(k string, m T) (float64, error)             { return Default.ZscoreF64(k, m) }

// agg is min,max, or others means sum
// when len(w) <> len(keys), not set WEIGHTS
func (p *Redis) Zunionstore(d string, keys []string, agg string, w []float64) (int, error) {
	args := Args(d, len(keys), keys)
	if agg == "MIN" || agg == "MAX" {
		args.Add("AGGREGATE", agg)
	}
	if len(w) > 0 && len(keys) == len(w) {
		args.Add("WEIGHTS").AddFlat(w)
	}
	return redis.Int(p.Do("ZUNIONSTORE", args...))
}
func Zunionstore(d string, keys []string, agg string, w []float64) (int, error) {
	return Default.Zunionstore(d, keys, agg, w)
}
