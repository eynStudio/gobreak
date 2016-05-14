package http

import (
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/eynstudio/gobreak"
)

func PostJson(url, post string) ([]byte, error) {
	resp, err := http.Post(url, "application/json", strings.NewReader(post))
	return GetRespBytes(resp, err)
}

//Post json with Header
func PostJsonWiHeader(url, post string, m M) ([]byte, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(post))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, _ := range m {
		req.Header.Set(k, m.GetStr(k))
	}
	return GetRespBytes(http.DefaultClient.Do(req))
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
