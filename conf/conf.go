package conf

import (
	"encoding/json"
	"io/ioutil"
)

func LoadJsonCfg(cfg interface{}, file string) (err error) {
	var content []byte

	if content, err = ioutil.ReadFile(file); err != nil {
		return err
	}

	return json.Unmarshal(content, cfg)
}
