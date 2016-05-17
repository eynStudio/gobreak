package gobreak

import (
	"bufio"
	"strconv"
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

func Str2F64(str string, defaultVal float64) float64 {
	if f, err := strconv.ParseFloat(str, 64); err == nil {
		return f
	}
	return defaultVal
}
