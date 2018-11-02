package jio

import (
	"reflect"
	"strings"
)

// NewContext Generates a context object with the provided data.
func NewContext(data interface{}) *Context {
	return &Context{
		root:   data,
		Value:  data,
		fields: make([]string, 0, 3),
	}
}

// Context contains data and toolkit
type Context struct {
	Value     interface{}
	Err       error
	root      interface{}
	fields    []string
	storage   map[string]interface{}
	skip      bool
	kindCache map[*interface{}]reflect.Kind
}

// Ref return the reference value.
// The reference path support use `.` access object property, just like javascript.
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

// FieldPath the field path of the current value.
func (ctx *Context) FieldPath() string {
	return strings.Join(ctx.fields, ".")
}

// Abort throw an error and skip the following check rules.
func (ctx *Context) Abort(err error) {
	ctx.Err = err
	ctx.skip = true
}

// Skip the following check rules.
func (ctx *Context) Skip() {
	ctx.skip = true
}

// Set is used to store a new key/value pair exclusively for this context.
func (ctx *Context) Set(name string, value interface{}) {
	if ctx.storage == nil {
		ctx.storage = make(map[string]interface{})
	}
	ctx.storage[name] = value
}

// Get returns the value for the given key, ie: (value, true).
func (ctx *Context) Get(name string) (interface{}, bool) {
	if ctx.storage == nil {
		ctx.storage = make(map[string]interface{})
	}
	value, ok := ctx.storage[name]
	return value, ok
}

// AssertKind assert the value type and cache.
func (ctx *Context) AssertKind(kind reflect.Kind) bool {
	if ctx.kindCache == nil {
		ctx.kindCache = make(map[*interface{}]reflect.Kind)
	}
	cachedKind, ok := ctx.kindCache[&ctx.Value]
	if !ok {
		cachedKind = reflect.TypeOf(ctx.Value).Kind()
		ctx.kindCache[&ctx.Value] = cachedKind
	}
	return cachedKind == kind
}
