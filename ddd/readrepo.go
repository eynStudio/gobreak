package ddd

import (
	"errors"
)

var (
	ErrCouldNotSaveModel = errors.New("could not save model")
	ErrModelNotFound     = errors.New("could not find model")
)

type ReadRepository interface {
	Save(GUID, interface{}) error
	Find(GUID) (interface{}, error)
	FindAll() ([]interface{}, error)
	Remove(GUID) error
}

type MemoryReadRepository struct {
	data map[GUID]interface{}
}

func NewMemoryReadRepository() *MemoryReadRepository {
	r := &MemoryReadRepository{
		data: make(map[GUID]interface{}),
	}
	return r
}

func (r *MemoryReadRepository) Save(id GUID, model interface{}) error {
	r.data[id] = model
	return nil
}

func (r *MemoryReadRepository) Find(id GUID) (interface{}, error) {
	if model, ok := r.data[id]; ok {
		return model, nil
	}

	return nil, ErrModelNotFound
}

func (r *MemoryReadRepository) FindAll() ([]interface{}, error) {
	models := []interface{}{}
	for _, model := range r.data {
		models = append(models, model)
	}
	return models, nil
}

func (r *MemoryReadRepository) Remove(id GUID) error {
	if _, ok := r.data[id]; ok {
		delete(r.data, id)
		return nil
	}

	return ErrModelNotFound
}
