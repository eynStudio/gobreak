package {{.Pkg}}

import (
	. "github.com/eynstudio/gobreak"
	"github.com/eynstudio/gobreak/db/mgo"
)

type {{.Name}} struct {
	Id     GUID   `bson:"_id,omitempty"`
	UserId GUID   `UserId`
}

func New{{.Name}}(uid GUID) *{{.Name}}{
	return &{{.Name}}{
		Id:mgo.NewGuid(),
		UserId:uid,
	}
}