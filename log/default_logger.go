package log

import(
	"log"
)

type DefaultLogger struct{}

func (p *DefaultLogger) Log(level int,msg string){
	log.Printf("[%d] %s\n",level, msg)
}

func UseDefaultLogger(){
	AddLogger(&DefaultLogger{})
}