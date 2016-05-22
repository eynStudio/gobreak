package gobreak

var OkStatus Status = Status{Code: 0, Msg: "OK"}

type Status struct {
	Code int
	Msg  string
}

func (p *Status) ErrMsg(msg string) *Status { return p.SetStatus(1, msg) }
func (p *Status) OkMsg(msg string) *Status  { return p.SetStatus(0, msg) }
func (p *Status) Ok() *Status               { return p.SetStatus(0, "OK") }
func (p *Status) IsOk() bool                { return p.Code == 0 }
func (p *Status) IsErr() bool               { return p.Code != 0 }

func (p *Status) SetStatus(code int, msg string) *Status {
	p.Code = code
	p.Msg = msg
	return p
}
