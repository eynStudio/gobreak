package captcha

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomDigits(length int) (b []byte) {
	b = make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = byte(rand.Intn(10))
	}
	return
}
