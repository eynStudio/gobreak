package gobreak

import (
	"encoding/json"
	"io/ioutil"
)

func LoadJson(file string, model T) bool {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return false
	}

	if err = json.Unmarshal(content, &model); err != nil {
		return false
	}

	return true
}
