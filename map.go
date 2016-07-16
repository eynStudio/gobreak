package gobreak

import (
	"log"
	"sort"
	"strconv"
)

type M map[string]T

func NewM() M { return make(M, 0) }

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

func (p M) GetGuid(k string) GUID { return GUID(p.GetStr(k)) }

func (p M) GetKeys() (keys []string) {
	for k := range p {
		keys = append(keys, k)
	}
	return
}

func (p M) GetSortedKeys() (keys []string) {
	keys = p.GetKeys()
	sort.Strings(keys)
	return
}

func getValAsInt(v T) int {
	switch t := v.(type) {
	case int:
		return t
	case float64:
		return int(t)
	case string:
		if i, err := strconv.Atoi(t); err == nil {
			return i
		}
		return 0
	default:
		log.Printf("%v is %v\n", v, t)
		return 0
	}
}

func getValAsF64(v T) float64 {
	switch t := v.(type) {
	case int:
		return float64(t)
	case float64:
		return t
	default:
		log.Printf("%v is %v\n", v, t)
		return 0
	}
}
