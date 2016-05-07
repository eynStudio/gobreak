package gobreak

var OkStatus MsgStatus = MsgStatus{Status: "ok"}

type MsgStatus struct {
	Msg    string
	Status string
}

func (p *MsgStatus) ErrMsg(msg string) *MsgStatus { return p.SetMsg("err", msg) }
func (p *MsgStatus) OkMsg(msg string) *MsgStatus  { return p.SetMsg("ok", msg) }
func (p *MsgStatus) Ok() *MsgStatus               { return p.SetMsg("ok", "") }

func (p *MsgStatus) SetMsg(status, msg string) *MsgStatus {
	p.Status = status
	p.Msg = msg
	return p
}
