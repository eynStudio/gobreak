package http

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	. "github.com/eynstudio/gobreak"
	"github.com/eynstudio/gox/encoding/jsonx"
)

var (
	NotSupportErr = errors.New("Data type not support")
	jsonHeader    = M{"Content-Type": "application/json"}
)

type Http struct {
	Error
	req    *http.Request
	resp   *http.Response
	reader io.Reader
}

func New() *Http { return &Http{} }

func (p *Http) Get(url string) *Http {
	p.NoErrExec(func() { p.req, p.Err = http.NewRequest("GET", url, nil) })
	return p
}

func (p *Http) Post(url string, data T, fmtType string) *Http {
	p.getReader(data, fmtType)
	p.NoErrExec(func() { p.req, p.Err = http.NewRequest("POST", url, p.reader) })
	return p
}

func (p *Http) Header(header M) *Http {
	if p.NotErr() {
		for k, _ := range header {
			p.req.Header.Set(k, header.GetStr(k))
		}
	}
	return p
}

func (p *Http) GetBytes() (data []byte) {
	p.NoErrExec(func() { data, p.Err = GetRespBytes(http.DefaultClient.Do(p.req)) })
	return data
}
func (p *Http) GetStr() (str string) { return string(p.GetBytes()) }
func (p *Http) GetJson(m T) *Http {
	data := p.GetBytes()
	log.Println(string(data))
	p.NoErrExec(func() { p.Err = json.Unmarshal(data, m) })
	return p
}
func (p *Http) GetXml(m T) *Http {
	data := p.GetBytes()
	p.NoErrExec(func() { p.Err = xml.Unmarshal(data, m) })
	return p
}

func (p *Http) Save(f string) *Http {
	p.resp, p.Err = http.DefaultClient.Do(p.req)
	defer p.resp.Body.Close()
	var file *os.File
	p.NoErrExec(func() { file, p.Err = os.Create(f) })
	p.NoErrExec(func() {
		_, p.Err = io.Copy(file, p.resp.Body)
		file.Close()
	})
	return p
}
func (p *Http) getReader(data T, objType string) {
	switch v := data.(type) {
	case io.Reader:
		p.reader = v
	case []byte:
		p.reader = bytes.NewReader(v)
	case string:
		p.reader = strings.NewReader(v)
	default:
		switch objType {
		case "json":
			p.reader = jsonx.Encode(data)
			//			var v []byte
			//			if v, p.Err = json.Marshal(data); p.NotErr() {
			//				p.reader = bytes.NewReader(v)
			//			}
		default:
			p.Err = NotSupportErr
		}
	}
}

func Get(url string) *Http                           { return New().Get(url) }
func GetWiHeader(url string, m M) *Http              { return New().Get(url).Header(m) }
func Post(url string, post T, fmtType string) *Http  { return New().Post(url, post, fmtType) }
func PostJson(url string, post T) *Http              { return New().Post(url, post, "json").Header(jsonHeader) }
func PostJsonWiHeader(url string, post T, m M) *Http { return PostJson(url, post).Header(m) }

func GetRespBytes(resp *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func GetReqIp(r *http.Request) string {
	if ipProxy := r.Header.Get("X-FORWARDED-FOR"); len(ipProxy) > 0 {
		return ipProxy
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
