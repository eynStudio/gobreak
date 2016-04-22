package store

import (
	"errors"
	"log"
	"reflect"

	. "github.com/eynstudio/gobreak"
	. "github.com/eynstudio/gobreak/db"
	. "github.com/eynstudio/gobreak/dddd/ddd"
	"github.com/eynstudio/gobreak/di"
)

var (
	ErrAggAlreadyRegistered = errors.New("aggregate is already registered")
	ErrAggNotRegistered     = errors.New("aggregate is not registered")
	aggMap                  = map[reflect.Type]reflect.Type{}
)

func Save(agg Agg) error {
	return nil
}

func Load(aggType reflect.Type, id GUID) (Agg, error) {
	if repoType, ok := aggMap[aggType]; ok {
		log.Println(repoType)
		if repo := di.Get(repoType); repo.IsValid() {
			r := repo.Interface().(Repo)
			a := r.Get(id)
			log.Println(a)

			agg := reflect.New(aggType).Interface().(Agg)
			log.Println(aggType, id, agg)
			return agg, nil
		}
	}

	//	if f, ok := p.callbacks[aggregateType]; ok {
	//		return p.EventStore.Load(f(id))
	//	} else {
	//		return nil, ErrAggregateNotRegistered
	//	}

	return nil, errors.New("error")
}

//func Load(agg Agg) (Agg, error) {
//	//	if name, ok := aggMap[reflect.TypeOf(agg)]; ok {
//	//		if repo, ok := repoMap[name]; ok {

//	//		}
//	//	}

//	return nil, nil
//}

func RegRepo(agg Agg, repo Repo) {
	aggMap[di.Ptr2Elem(reflect.TypeOf(agg))] = reflect.TypeOf(repo)
}
