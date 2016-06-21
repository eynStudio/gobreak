package http

import (
	"strings"
	"log"
	"testing"
)

func Test_getReader(t *testing.T) {
	log.Println("hi")
	
//	hi:=Hi{}
	rr:=strings.NewReader("ssss")
	r,err:=getReader(rr,"")
	
		log.Println(r,err)

}

type Hi struct{
	Hi string
}