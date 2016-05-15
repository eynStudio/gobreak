package redis

import (
	"encoding/json"

	. "github.com/eynstudio/gobreak"

	"github.com/garyburd/redigo/redis"
)

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

func (p *Cmd) Bool() (b bool) {
	b, p.Err = redis.Bool(p.reply, p.Err)
	return
}

func (p *Cmd) String() (s string) {
	s, p.Err = redis.String(p.reply, p.Err)
	return
}

func (p *Cmd) Values() (t []interface{}) {
	t, p.Err = redis.Values(p.reply, p.Err)
	return
}

func (p *Cmd) As(m T) *Cmd {
	if vals := p.Values(); p.NotErr() {
		if len(vals) == 0 {
			p.Err = redis.ErrNil
			return p
		}
		p.Err = redis.ScanStruct(vals, m)
	}
	return p
}

func (p *Cmd) Json(m T) *Cmd {
	p.NoErrExec(func() { p.Err = json.Unmarshal(p.Bytes(), &m) })
	return p
}

func (p Cmd) IsNotFound() bool { return p.Error.IsErr() && p.Err == redis.ErrNil }
func (p Cmd) IsErr() bool      { return p.Error.IsErr() && p.Err != redis.ErrNil }
