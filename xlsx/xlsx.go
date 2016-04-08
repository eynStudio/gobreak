package xlsx

import (
	. "github.com/eynstudio/gobreak"
	"github.com/tealeg/xlsx"
)

func Open(fname string) (*xlsx.File, error) {
	return xlsx.OpenFile(fname)
}

func MustOpen(fname string) *xlsx.File {
	xlFile, err := Open(fname)
	Must(err)
	return xlFile
}
