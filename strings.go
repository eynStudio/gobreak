package gobreak

import (
	"bufio"
	"strings"
)

func StrSplit(str string) (lines []string) {
	s := strings.NewReader(str)
	bs := bufio.NewScanner(s)
	for bs.Scan() {
		l := bs.Text()
		lines = append(lines, l)
	}
	return
}
