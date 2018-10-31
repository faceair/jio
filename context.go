package jio

import (
	"reflect"
	"strings"
)

func NewContext(data interface{}) *Context {
	return &Context{
		root:      data,
		Value:     data,
		fields:    make([]string, 0, 3),
		storage:   make(map[string]interface{}),
		kindCache: make(map[*interface{}]reflect.Kind),
	}
}

type Context struct {
	Value     interface{}
	root      interface{}
	fields    []string
	storage   map[string]interface{}
	err       error
	skip      bool
	kindCache map[*interface{}]reflect.Kind
}

func (ctx *Context) Ref(refPath string) (value interface{}, ok bool) {
	fields := strings.Split(refPath, ".")
	value = ctx.root
	var valueMap map[string]interface{}
	for _, field := range fields {
		valueMap, ok = value.(map[string]interface{})
		if !ok {
			return
		}
		value, ok = valueMap[field]
		if !ok {
			return
		}
	}
	return
}

func (ctx *Context) FieldPath() string {
	return strings.Join(ctx.fields, ".")
}

func (ctx *Context) Abort(err error) {
	ctx.err = err
	ctx.skip = true
}

func (ctx *Context) Skip() {
	ctx.skip = true
}

func (ctx *Context) Set(name string, value interface{}) {
	ctx.storage[name] = value
}

func (ctx *Context) Get(name string) (interface{}, bool) {
	value, ok := ctx.storage[name]
	return value, ok
}

func (ctx *Context) AssertKind(kind reflect.Kind) bool {
	cachedKind, ok := ctx.kindCache[&ctx.Value]
	if !ok {
		cachedKind = reflect.TypeOf(ctx.Value).Kind()
		ctx.kindCache[&ctx.Value] = cachedKind
	}
	return cachedKind == kind
}
