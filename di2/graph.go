package di2

import (
	"fmt"
	"reflect"

	"github.com/eynstudio/gobreak/log"
)

type Graph struct {
	Logger      log.ILogger
	unnamed     []*Obj
	unnamedType map[reflect.Type]bool
	named       map[string]*Obj
}

func (p *Graph) RegAs(name string, val interface{}) error {
	if err := p.reg(&Obj{val: val, name: name}); err != nil {
		return err
	}
	return p.Populate()
}

func (p *Graph) Reg(values ...interface{}) error {
	for _, v := range values {
		if err := p.reg(&Obj{val: v}); err != nil {
			return err
		}
	}
	return p.Populate()
}

func (p *Graph) reg(objs ...*Obj) error {
	p.checkInit()
	for _, o := range objs {
		o.diType = reflect.TypeOf(o.val)
		o.diVal = reflect.ValueOf(o.val)
		if o.name == "" {
			if !IsStructPtr(o.diType) {
				return fmt.Errorf("unnamed obj must be a pointer to struct, %s, %v", o.diType, o.val)
			}
			if !o.private {
				if p.unnamedType[o.diType] {
					return fmt.Errorf("two unnamed instances of type *%s.%s", o.diType.Elem().PkgPath(), o.diType.Elem().Name())
				}
				p.unnamedType[o.diType] = true
			}
			p.unnamed = append(p.unnamed, o)
		} else {
			if p.named[o.name] != nil {
				return fmt.Errorf("two instances named %s", o.name)
			}
			p.named[o.name] = o
		}
		if p.Logger != nil {
			if o.created {
				p.logf("created %s", o)
			} else if o.embedded {
				p.logf("provided embedded %s", o)
			} else {
				p.logf("provided %s", o)
			}
		}
	}
	return nil
}

func (p *Graph) Populate() error {
	for _, o := range p.named {
		if err := p.populateExplicit(o); err != nil {
			return err
		}
	}

	i := 0
	for {
		if i == len(p.unnamed) {
			break
		}
		o := p.unnamed[i]
		i++
		if err := p.populateExplicit(o); err != nil {
			return err
		}
	}

	for _, o := range p.unnamed {
		if err := p.populateUnnamedInterface(o); err != nil {
			return err
		}
	}

	for _, o := range p.named {
		if err := p.populateUnnamedInterface(o); err != nil {
			return err
		}
	}
	return nil
}

func (p *Graph) populateExplicit(o *Obj) error {
	if o.name != "" && !IsStructPtr(o.diType) {
		return nil
	}

StructLoop:
	for i := 0; i < o.diVal.Elem().NumField(); i++ {
		field := o.diVal.Elem().Field(i)
		fieldType := field.Type()
		fieldTag := o.diType.Elem().Field(i).Tag
		fieldName := o.diType.Elem().Field(i).Name
		tag := parseTag(fieldTag)
		if tag == nil {
			continue
		}
		if !field.CanSet() {
			continue
		}
		if tag.Inline && fieldType.Kind() != reflect.Struct {
			return fmt.Errorf("inline requested on non inlined field %s in type %s", fieldName, o.diType)
		}
		if !IsNilOrZero(field, fieldType) {
			continue
		}

		if tag.Name != "*" {
			existing := p.named[tag.Name]
			if existing == nil {
				return fmt.Errorf("did not find object named %s required by field %s in type %s", tag.Name, fieldName, o.diType)
			}

			if !existing.diType.AssignableTo(fieldType) {
				return fmt.Errorf(
					"object named %s of type %s is not assignable to field %s (%s) in type %s",
					tag.Name, fieldType, fieldName, existing.diType, o.diType)
			}

			field.Set(reflect.ValueOf(existing.val))
			p.logf("assigned %s to field %s in %s", existing, fieldName, o)
			o.setField(fieldName, existing)
			continue StructLoop
		}

		if fieldType.Kind() == reflect.Struct {
			if tag.Private {
				return fmt.Errorf("cannot use private inject on inline struct on field %s in type %s", fieldName, o.diType)
			}
			if !tag.Inline {
				return fmt.Errorf("inline struct on field %s in type %s requires an explicit \"inline\" tag", fieldName, o.diType)
			}

			err := p.reg(&Obj{
				val:      field.Addr().Interface(),
				private:  true,
				embedded: o.diType.Elem().Field(i).Anonymous,
			})
			if err != nil {
				return err
			}
			continue
		}

		if fieldType.Kind() == reflect.Interface {
			continue
		}

		if fieldType.Kind() == reflect.Map {
			if !tag.Private {
				return fmt.Errorf("inject on map field %s in type %s must be named or private", fieldName, o.diType)
			}

			field.Set(reflect.MakeMap(fieldType))
			p.logf("made map for field %s in %s", fieldName, o)
			continue
		}

		if !IsStructPtr(fieldType) {
			return fmt.Errorf("found inject tag on unsupported field %s in type %s", fieldName, o.diType)
		}

		if !tag.Private {
			for _, existing := range p.unnamed {
				if existing.private {
					continue
				}
				if existing.diType.AssignableTo(fieldType) {
					field.Set(reflect.ValueOf(existing.val))
					p.logf("assigned existing %s to field %s in %s", existing, fieldName, o)
					o.setField(fieldName, existing)
					continue StructLoop
				}
			}
		}

		newValue := reflect.New(fieldType.Elem())
		newObject := &Obj{val: newValue.Interface(), private: tag.Private, created: true}

		if err := p.reg(newObject); err != nil {
			return err
		}

		field.Set(newValue)
		p.logf("assigned newly created %s to field %s in %s", newObject, fieldName, o)
		o.setField(fieldName, newObject)
	}

	return nil
}

func (p *Graph) populateUnnamedInterface(o *Obj) error {
	if o.name != "" && !IsStructPtr(o.diType) {
		return nil
	}

	for i := 0; i < o.diVal.Elem().NumField(); i++ {
		field := o.diVal.Elem().Field(i)
		fieldType := field.Type()
		fieldTag := o.diType.Elem().Field(i).Tag
		fieldName := o.diType.Elem().Field(i).Name
		tag := parseTag(fieldTag)

		if tag == nil {
			continue
		}
		if fieldType.Kind() != reflect.Interface {
			continue
		}
		if tag.Private {
			return fmt.Errorf("found private inject tag on interface field %s in type %s", fieldName, o.diType)
		}

		if !IsNilOrZero(field, fieldType) {
			continue
		}

		if tag.Name != "*" {
			panic(fmt.Sprintf("unhandled named instance with name %s", tag.Name))
		}

		var found *Obj
		for _, existing := range p.unnamed {
			if existing.private {
				continue
			}
			if existing.diType.AssignableTo(fieldType) {
				if found != nil {
					return fmt.Errorf(
						"found two assignable values for field %s in type %s. one type "+
							"%s with value %v and another type %s with value %v",
						fieldName, o.diType, found.diType, found.val, existing.diType, existing.diVal,
					)
				}
				found = existing
				field.Set(reflect.ValueOf(existing.val))
				p.logf("assigned existing %s to interface field %s in %s", existing, fieldName, o)
				o.setField(fieldName, existing)
			}
		}

		if found == nil {
			return fmt.Errorf("found no assignable value for field %s in type %s", fieldName, o.diType)
		}
	}
	return nil
}

func (p *Graph) logf(format string, v ...interface{}) {
	if p.Logger != nil {
		p.Logger.Printf(format, v...)
	}
}

func (p *Graph) checkInit() {
	if p.unnamedType == nil {
		p.unnamedType = make(map[reflect.Type]bool)
	}
	if p.named == nil {
		p.named = make(map[string]*Obj)
	}
}
