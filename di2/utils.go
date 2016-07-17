package di2

import "reflect"

func IsStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

func IsNilOrZero(v reflect.Value, t reflect.Type) bool {
	switch v.Kind() {
	default:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(t).Interface())
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
}

var (
	injectOnly    = &tag{}
	injectPrivate = &tag{Name: "", Private: true}
	injectInline  = &tag{Name: "", Inline: true}
)

type tag struct {
	Name    string
	Inline  bool
	Private bool
}

func parseTag(tags reflect.StructTag) *tag {
	val := tags.Get("di")
	switch val {
	case "":
		return nil
	case "inline":
		return injectInline
	case "!":
		return injectPrivate
	case "*":
		return injectOnly
	default:
		return &tag{Name: val}
	}

}
