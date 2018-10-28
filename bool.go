package jio

import "fmt"

func Bool() *BoolSchema {
	return &BoolSchema{
		rules: make([]func(*Context), 0, 3),
	}
}

var _ Schema = new(BoolSchema)

type BoolSchema struct {
	required *bool
	rules    []func(*Context)
}

func (b *BoolSchema) Required() *BoolSchema {
	b.required = boolPtr(true)
	b.rules = append([]func(*Context){func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Abort(fmt.Errorf("field `%s` is required", ctx.FieldPath()))
		}
	}}, b.rules...)
	return b
}

func (b *BoolSchema) Optional() *BoolSchema {
	b.required = boolPtr(false)
	b.rules = append([]func(*Context){func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Skip()
		}
	}}, b.rules...)
	return b
}

func (b *BoolSchema) Default(value bool) *BoolSchema {
	b.required = boolPtr(false)
	b.rules = append([]func(*Context){func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Value = value
		}
	}}, b.rules...)
	return b
}

func (b *BoolSchema) Truthy(values ...interface{}) *BoolSchema {
	b.rules = append(b.rules, func(ctx *Context) {
		for _, v := range values {
			if v == ctx.Value {
				ctx.Value = true
			}
		}
	})
	return b
}

func (b *BoolSchema) Falsy(values ...interface{}) *BoolSchema {
	b.rules = append(b.rules, func(ctx *Context) {
		for _, v := range values {
			if v == ctx.Value {
				ctx.Value = false
			}
		}
	})
	return b
}

func (b *BoolSchema) Validate(ctx *Context) {
	if b.required == nil {
		b.Optional()
	}
	for _, rule := range b.rules {
		rule(ctx)
		if ctx.skip {
			return
		}
	}
	if ctx.err == nil {
		if _, ok := (ctx.Value).(bool); !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not boolean", ctx.FieldPath(), ctx.Value))
		}
	}
}
