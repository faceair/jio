package jio

import (
	"fmt"
)

var _ Schema = new(AnySchema)

func Any() *AnySchema {
	return &AnySchema{
		rules: make([]func(*Context), 0, 3),
	}
}

type AnySchema struct {
	baseSchema

	required *bool
	rules    []func(*Context)
}

func (a *AnySchema) SetPriority(priority int) *AnySchema {
	a.priority = priority
	return a
}

func (a *AnySchema) PrependTransform(f func(*Context)) *AnySchema {
	a.rules = append([]func(*Context){f}, a.rules...)
	return a
}

func (a *AnySchema) Transform(f func(*Context)) *AnySchema {
	a.rules = append(a.rules, f)
	return a
}

func (a *AnySchema) Required() *AnySchema {
	a.required = boolPtr(true)
	return a.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Abort(fmt.Errorf("field `%s` is required", ctx.FieldPath()))
		}
	})
}

func (a *AnySchema) Optional() *AnySchema {
	a.required = boolPtr(false)
	return a.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Skip()
		}
	})
}

func (a *AnySchema) Default(value interface{}) *AnySchema {
	a.required = boolPtr(false)
	return a.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Value = value
		}
	})
}

func (a *AnySchema) When(refPath string, condition interface{}, then Schema) *AnySchema {
	return a.Transform(func(ctx *Context) { a.when(ctx, refPath, condition, then) })
}

func (a *AnySchema) Valid(values ...interface{}) *AnySchema {
	return a.Transform(func(ctx *Context) {
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
}

func (a *AnySchema) Equal(value interface{}) *AnySchema {
	return a.Transform(func(ctx *Context) {
		if value != ctx.Value {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not %v", ctx.FieldPath(), ctx.Value, value))
			return
		}
	})
}

func (a *AnySchema) Validate(ctx *Context) {
	if a.required == nil {
		a.Optional()
	}
	for _, rule := range a.rules {
		rule(ctx)
		if ctx.skip {
			return
		}
	}
}
