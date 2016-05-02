package store

import (
	"errors"
	"log"
	"reflect"

	. "github.com/eynstudio/gobreak"
	. "github.com/eynstudio/gobreak/db"
	. "github.com/eynstudio/gobreak/dddd/ddd"
	"github.com/eynstudio/gobreak/dddd/eventbus"
	"github.com/eynstudio/gobreak/di"
)

var (
	ErrAggAlreadyRegistered = errors.New("aggregate is already registered")
	ErrAggNotRegistered     = errors.New("aggregate is not registered")
	aggMap                  = map[reflect.Type]reflect.Type{}
)

func Save(agg Agg) error {
	events := agg.GetUncommittedEvents()
	if len(events) == 0 {
		return nil
	}

	for _, event := range events {
		eventbus.Publish(event)
	}

	agg.ClearUncommittedEvents()

	if repo, ok := getRepoByAgg(reflect.TypeOf(agg)); ok {
		if agg.IsDeleted() {
			repo.Del(agg.ID())
		} else {
			repo.Save(agg.ID(), agg.Root())
		}
	} else {
		log.Println("no repo")
	}
	return nil
}

func Load(aggType reflect.Type, id GUID) (Agg, error) {
	if repoType, ok := aggMap[di.Ptr2Elem(aggType)]; ok {
		if repo := di.Get(repoType); repo.IsValid() {
			r := repo.Interface().(Repo)
			agg := reflect.New(aggType).Interface().(Agg)
			r.GetAs(id, agg.Root())
			return agg, nil
		}
	}
	return nil, errors.New("error")
}

func getRepoByAgg(aggType reflect.Type) (Repo, bool) {
	if repoType, ok := aggMap[di.Ptr2Elem(aggType)]; ok {
		if repo := di.Get(repoType); repo.IsValid() {
			if r := repo.Interface().(Repo); r != nil {
				return r, true
			}
		}
	}
	return nil, false
}

func RegRepo(agg Agg, repo Repo) {
	aggMap[di.Ptr2Elem(reflect.TypeOf(agg))] = reflect.TypeOf(repo)
}
