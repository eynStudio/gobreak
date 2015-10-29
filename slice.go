package gobreak

import (
	"reflect"
)

type slice struct {
	val reflect.Value
}

func Slice(ptr T) *slice {
	MustSlicePtr(ptr)
	return &slice{reflect.ValueOf(ptr).Elem()}
}

func (s *slice) Find(t T) int {
	for i := 0; i < s.val.Len(); i++ {
		if reflect.DeepEqual(s.val.Index(i).Interface(), t) {
			return i
		}
	}
	return -1
}

func (s *slice) FindEntity(id GUID) Entity {
	for i := 0; i < s.val.Len(); i++ {
		entity := s.val.Index(i).Interface().(Entity)
		if entity.ID() == id {
			return entity
		}
	}
	return nil
}
func (s *slice) FindEntityIndex(id GUID) int {
	for i := 0; i < s.val.Len(); i++ {
		if s.val.Index(i).Interface().(Entity).ID() == id {
			return i
		}
	}
	return -1
}
func (s *slice) ReplaceEntity(entity Entity) *slice {
	val_entity := reflect.ValueOf(entity)
	for i := 0; i < s.val.Len(); i++ {
		if s.val.Index(i).Interface().(Entity).ID() == entity.ID() {
			s.val.Index(i).Set(val_entity)
			return s
		}
	}
	s.val.Set(reflect.Append(s.val.Slice(0, s.val.Len()), val_entity))
	return s
}

func (s *slice) FindBy(f func(T) bool) int {
	for i := 0; i < s.val.Len(); i++ {
		if f(s.val.Index(i).Interface()) {
			return i
		}
	}
	return -1
}

func (s *slice) Each(f func(T, int)) *slice {
	for i := 0; i < s.val.Len(); i++ {
		f(s.val.Index(i).Interface(), i)
	}
	return s
}

func (s *slice) RemoveAt(i int) *slice {
	if i >= 0 && i < s.val.Len() {
		s.val.Set(reflect.AppendSlice(s.val.Slice(0, i), s.val.Slice(i+1, s.val.Len())))
	}
	return s
}

func (s *slice) Remove(t T) *slice              { return s.RemoveAt(s.Find(t)) }
func (s *slice) RemoveEntity(id GUID) *slice    { return s.RemoveAt(s.FindEntityIndex(id)) }
func (s *slice) RemoveBy(f func(T) bool) *slice { return s.RemoveAt(s.FindBy(f)) }
