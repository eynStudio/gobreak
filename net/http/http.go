package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	. "github.com/eynstudio/gobreak"
)

var NotSupportErr = errors.New("Data type not support")

func PostWiHeader(url string, post T, header M) ([]byte, error) {
	r, err := getReader(post, "")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		return nil, err
	}

	for k, _ := range header {
		req.Header.Set(k, header.GetStr(k))
	}
	return GetRespBytes(http.DefaultClient.Do(req))
}

func PostJson(url string, post T) ([]byte, error) { return PostJsonWiHeader(url, post, nil) }

func PostJsonWiHeader(url string, post T, m M) ([]byte, error) {
	r, err := getReader(post, "json")
	if err != nil {
		return nil, err
	}
	if m == nil {
		m = M{}
	}
	m["Content-Type"] = "application/json"
	return PostWiHeader(url, r, m)
}

func Get(url string) ([]byte, error) {
	return GetRespBytes(http.Get(url))
}

//Get with Header
func GetWiHeader(url string, m M) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, _ := range m {
		req.Header.Set(k, m.GetStr(k))
	}
	return GetRespBytes(http.DefaultClient.Do(req))
}

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

func getReader(data T, objType string) (r io.Reader, err error) {
	switch v := data.(type) {
	case io.Reader:
		return v, nil
	case []byte:
		return bytes.NewReader(v), nil
	case string:
		return strings.NewReader(v), nil
	default:
		switch objType {
		case "json":
			if v, err := json.Marshal(data); err != nil {
				return nil, err
			} else {
				return bytes.NewReader(v), nil
			}
		default:
			return nil, NotSupportErr
		}
	}
}
