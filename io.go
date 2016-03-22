package gobreak

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
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

func ReadLines(f string) (lines []string, err error) {
	var file *os.File
	if file, err = os.Open(f); err != nil {
		return
	}

	reader := bufio.NewReader(file)
	bs := bufio.NewScanner(reader)
	for bs.Scan() {
		l := bs.Text()
		lines = append(lines, l)
	}
	return
}
