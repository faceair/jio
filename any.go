package jio

import (
	"fmt"
)

type Schema interface {
	Validate(*Context)
}

var _ Schema = new(AnySchema)

func Any() *AnySchema {
	return &AnySchema{
		rules: make([]func(*Context), 0, 3),
	}
}

type AnySchema struct {
	rules []func(*Context)
}

func (a *AnySchema) Required() *AnySchema {
	a.rules = append(a.rules, func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Abort(fmt.Errorf("field `%s` is required", ctx.FieldPath()))
		}
	})
	return a
}

func (a *AnySchema) Optional() *AnySchema {
	a.rules = append(a.rules, func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Skip()
		}
	})
	return a
}

func (a *AnySchema) Default(value interface{}) *AnySchema {
	a.rules = append(a.rules, func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Value = value
		}
	})
	return a
}

func (a *AnySchema) Valid(values ...interface{}) *AnySchema {
	a.rules = append(a.rules, func(ctx *Context) {
		var isValid bool
		for _, v := range values {
			if v == ctx.Value {
				isValid = true
				break
			}
		}
		if !isValid {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not in %v", ctx.FieldPath(), ctx.Value, values))
			return
		}
	})
	return a
}

func (a *AnySchema) Transform(f func(*Context)) *AnySchema {
	a.rules = append(a.rules, f)
	return a
}

func (a *AnySchema) Validate(ctx *Context) {
	for _, rule := range a.rules {
		rule(ctx)
		if ctx.skip {
			return
		}
	}
}
