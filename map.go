package gobreak

import (
	"log"
)

type M map[string]T

func (p M) HasKey(k string) bool {
	_, ok := p[k]
	return ok
}
func (p M) GetStrOr(k, or string) string {
	if v, ok := p[k]; ok {
		return v.(string)
	}
	return or
}
func (p M) GetStr(k string) string { return p.GetStrOr(k, "") }

func (p M) GetIntOr(k string, or int) int {
	if v, ok := p[k]; ok {
		return getValAsInt(v)
	}
	return or
}
func (p M) GetInt(k string) int { return p.GetIntOr(k, 0) }

func (p M) GetF64Or(k string, or float64) float64 {
	if v, ok := p[k]; ok {
		return getValAsF64(v)
	}
	return or
}
func (p M) GetF64(k string) float64 { return p.GetF64Or(k, 0) }

func (p M) GetBoolOr(k string, or bool) bool {
	if v, ok := p[k]; ok {
		return v.(bool)
	}
	return or
}
func (p M) GetBool(k string) bool { return p.GetBoolOr(k, false) }

func getValAsInt(v T) int {
	switch t := v.(type) {
	case int:
		return v.(int)
	case float64:
		return int(v.(float64))
	default:
		log.Println("%v is %v", v, t)
		return 0
	}
}

func getValAsF64(v T) float64 {
	switch t := v.(type) {
	case int:
		return float64(v.(int))
	case float64:
		return v.(float64)
	default:
		log.Println("%v is %v", v, t)
		return 0
	}
}
