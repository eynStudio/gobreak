package captcha

import (
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomDigits(length int) (str string) {
	for i := 0; i < length; i++ {
		str += strconv.Itoa(rand.Intn(10))
	}
	return
}
