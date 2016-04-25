package captcha

import (
	"log"
	"testing"
)

func Test_RandomDigits(t *testing.T) {
	for i := 0; i < 10; i++ {
		log.Println(RandomDigits(4))
	}
}
