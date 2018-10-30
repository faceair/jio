package jio

import (
	"strings"
)

func NewContext(data interface{}) *Context {
	return &Context{
		Root:    data,
		Value:   data,
		Fields:  make([]string, 0, 3),
		storage: make(map[string]interface{}),
	}
}

type Context struct {
	Root    interface{}
	Fields  []string
	Value   interface{}
	storage map[string]interface{}
	err     error
	skip    bool
}

func (ctx *Context) Ref(refPath string) (value interface{}, ok bool) {
	fields := strings.Split(refPath, ".")
	value = ctx.Root
	for _, field := range fields {
		value, ok = value.(map[string]interface{})[field]
		if !ok {
			return
		}
	}
	return
}

func (ctx *Context) FieldPath() string {
	return strings.Join(ctx.Fields, ".")
}

func (ctx *Context) Enter(fields []string, value interface{}) {
	ctx.Fields = fields
	ctx.Value = value
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
