package jio

import "strings"

func NewContext(data interface{}) *Context {
	return &Context{
		Fields: make([]string, 0, 3),
		Value:  data,
	}
}

type Context struct {
	Fields []string
	Value  interface{}
	err    error
	skip   bool
}

func (ctx *Context) FieldPath() string {
	return strings.Join(ctx.Fields, ".")
}

func (ctx *Context) Abort(err error) {
	ctx.err = err
	ctx.skip = true
}

func (ctx *Context) Skip() {
	ctx.skip = true
}
