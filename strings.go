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
func Strs2F64(strs []string, defaultVal float64) (f []float64) {
	for _, it := range strs {
		f = append(f, Str2F64(it, defaultVal))
	}
	return
}
func Str2Int(str string, defaultVal int) int {
	if f, err := strconv.Atoi(str); err == nil {
		return f
	}
	return defaultVal
}
func StrFromF64(val float64) string {
	return strconv.FormatFloat(val, 'f', -1, 64)
}
