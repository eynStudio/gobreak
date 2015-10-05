package converts

import (
	"strconv"
)

func StrParseF32(str string) (float32, error) {
	f, err := strconv.ParseFloat(str, 32)
	return float32(f), err
}

func Str2F32(str string, defaultVal float32) float32 {
	if f, err := strconv.ParseFloat(str, 32); err == nil {
		return float32(f)
	}
	return defaultVal
}
