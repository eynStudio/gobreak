package gobreak

var OkStatus = &Status{Code: 0, Msg: "OK"}

type IStatus interface {
	SetStatus(code int, msg string) IStatus
	ErrMsg(msg string) IStatus
	OkMsg(msg string) IStatus
	Ok() IStatus
	IsOk() bool
	IsErr() bool
	GetStatus() Status
}

type Status struct {
	Code int
	Msg  string
}

func NewStatus() *Status              { return &Status{} }
func NewErrStatus(msg string) *Status { return &Status{Code: 1, Msg: msg} }
func NewStatusErr(err error, ok_msg, err_msg string) IStatus {
	return NewStatus().OkErrMsg(err, ok_msg, err_msg)
}

func (p *Status) OkErrMsg(err error, ok_msg, err_msg string) IStatus {
	if err == nil {
		return p.OkMsg(ok_msg)
	}
	return p.ErrMsg(err_msg)
}

func (p *Status) ErrMsg(msg string) IStatus { return p.SetStatus(1, msg) }
func (p *Status) OkMsg(msg string) IStatus  { return p.SetStatus(0, msg) }
func (p *Status) Ok() IStatus               { return p.SetStatus(0, "OK") }
func (p *Status) IsOk() bool                { return p.Code == 0 }
func (p *Status) IsErr() bool               { return p.Code != 0 }

func (p *Status) SetStatus(code int, msg string) IStatus {
	p.Code = code
	p.Msg = msg
	return p
}
func (p *Status) GetStatus() Status {
	return *p
}
