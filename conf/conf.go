package conf

import (
	"encoding/json"
	"io/ioutil"
	"github.com/eynstudio/gobreak"
)

func LoadJsonCfg(cfg interface{}, file string) (err error) {
	var content []byte

	if content, err = ioutil.ReadFile(file); err != nil {
		return err
	}

	return json.Unmarshal(content, cfg)
}

func MustLoadJsonCfg(cfg interface{}, file string){
	gobreak.Must(LoadJsonCfg(cfg,file))
}