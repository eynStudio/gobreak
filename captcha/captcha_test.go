package captcha

import (
	"log"
	"testing"
)

func Test_RandomDigits(t *testing.T) {
	for i := 0; i < 10; i++ {
		code := RandomDigits(4)
		log.Println(string(code))
	}
}
