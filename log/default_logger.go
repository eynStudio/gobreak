package log

import(
	"log"
)

type DefaultLogger struct{}

func (p *DefaultLogger) Log(level LogLevel,msg string){
	log.Printf("[%s] %s\n",LogLevelName[level], msg)
}

func UseDefaultLogger(){
	AddLogger(&DefaultLogger{})
}