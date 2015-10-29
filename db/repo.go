package db

import (
	"errors"
	. "github.com/eynstudio/gobreak"
)

var (
	ErrModelNotFound = errors.New("could not find model")
)

type Repo interface {
	All() []T
	Get(id T) T
	Save(id T, m T)
	Del(id T)
}

type MemoryRepo struct {
	data map[GUID]T
}

func NewMemoryRepoRepo() *MemoryRepo {
	return &MemoryRepo{
		data: make(map[GUID]T),
	}
}

func (p *MemoryRepo) Save(id GUID, model T) error {
	p.data[id] = model
	return nil
}

func (p *MemoryRepo) Get(id GUID) T {
	if model, ok := p.data[id]; ok {
		return model
	}
	return nil
}

func (p *MemoryRepo) All() []T {
	models := []T{}
	for _, model := range p.data {
		models = append(models, model)
	}
	return models
}

func (p *MemoryRepo) Del(id GUID) error {
	if _, ok := p.data[id]; ok {
		delete(p.data, id)
		return nil
	}

	return ErrModelNotFound
}
