package datelog

import (
	"fmt"
	"os"
	"time"

	. "github.com/eynstudio/gobreak"
)

type DateLog struct {
	Error
	outPath string
	curDate time.Time
	file    *os.File
}

func New(path string) *DateLog {
	l := &DateLog{outPath: path, curDate: time.Now()}
	l.Err = os.MkdirAll(path, 0644)
	l.openFile()
	return l
}

func (p *DateLog) Log(v ...interface{}) {
	p.checkFile()
	fmt.Fprint(p.file, v...)
}

func (p *DateLog) Logf(format string, v ...interface{}) {
	p.checkFile()
	fmt.Fprintf(p.file, format, v...)
}

func (p *DateLog) Logln(v ...interface{}) {
	p.checkFile()
	fmt.Fprintln(p.file, v...)
}

func (p *DateLog) checkFile() {
	if p.curDate.Day() != time.Now().Day() {
		p.Close()
		p.curDate = time.Now()
		p.openFile()
	}
}

func (p *DateLog) getFileName() string {
	return fmt.Sprintf("%s/%s.log", p.outPath, p.curDate.Format("20060102"))
}

func (p *DateLog) openFile() {
	p.file, p.Err = os.OpenFile(p.getFileName(), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0664)
	p.LogErr()
}

func (p *DateLog) Close() {
	if p.file != nil {
		p.file.Close()
	}
}
