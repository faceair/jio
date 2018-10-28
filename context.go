package jio

import "strings"

type Context struct {
	Fields []string
	Value  interface{}
}

func (ctx *Context) FieldPath() string {
	return strings.Join(ctx.Fields, ".")
}

func (ctx *Context) Next() {

}
