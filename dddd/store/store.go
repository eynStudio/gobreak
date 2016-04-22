package store

import (
	"errors"
	. "github.com/eynstudio/gobreak/db"
	. "github.com/eynstudio/gobreak/dddd"
	"reflect"
)

var (
	ErrAggAlreadyRegistered = errors.New("aggregate is already registered")
	ErrAggNotRegistered     = errors.New("aggregate is not registered")
	aggMap                  = map[reflect.Type]string{}
	repoMap                 = map[string]Repo{}
)

func Save(agg Agg) error {
	return nil
}

func Load(agg Agg) (Agg, error) {
	//	if name, ok := aggMap[reflect.TypeOf(agg)]; ok {
	//		if repo, ok := repoMap[name]; ok {

	//		}
	//	}

	return nil, nil
}

func RegRepo(agg Agg, repo Repo) {
	repoMap[repo.GetName()] = repo
	aggMap[reflect.TypeOf(agg)] = repo.GetName()
}
