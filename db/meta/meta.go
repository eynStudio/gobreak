package meta

import (
//	. "github.com/eynstudio/gobreak"
)

type MetaDb struct {
	Mc   string
	Bz   string
	Tbls []MetaTbl
}

type MetaTbl struct {
	Mc   string
	Bz   string
	Cols []MetaCol
}

type MetaCol struct {
	Mc         string
	Lx         string
	Len        int
	Scale      int
	IsKey      bool
	Nullable   bool
	HasDefault bool
	Default    string
	Bz         string
}
