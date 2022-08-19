package symbolproc

import (
	"github.com/traefik/yaegi/interp"
	"path"
	"reflect"
	"regexp"
	"strings"
)

func Clean(symbols interp.Exports, prefix string) {
	for key, value := range symbols {
		delete(symbols, key)
		newKey := strings.TrimPrefix(key, prefix)
		symbols[newKey] = value
	}
}

func in(needle string, haystack []string) bool {
	for _, item := range haystack {
		r := regexp.MustCompile(item)
		if r.MatchString(needle) {
			return true
		}
	}
	return false
}

func AllowList(symbols interp.Exports, paths ...string) {
	for key := range symbols {
		if !in(path.Dir(key), paths) {
			delete(symbols, key)
		}
	}
}

func DenyList(symbols interp.Exports, paths ...string) {
	for key := range symbols {
		if in(path.Dir(key), paths) {
			delete(symbols, key)
		}
	}
}

func Pluck(values map[string]reflect.Value, items ...string) {
	for _, item := range items {
		delete(values, item)
	}
}
