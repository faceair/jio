package jio

import "fmt"

func Bool() *BoolSchema {
	return &BoolSchema{
		rules: make([]func(*Context), 0, 3),
	}
}

var _ Schema = new(BoolSchema)

type BoolSchema struct {
	rules []func(*Context)
}

func (b *BoolSchema) Required() *BoolSchema {
	b.rules = append(b.rules, func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Abort(fmt.Errorf("field `%s` is required", ctx.FieldPath()))
		}
	})
	return b
}

func (b *BoolSchema) Optional() *BoolSchema {
	b.rules = append(b.rules, func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Skip()
		}
	})
	return b
}

func (b *BoolSchema) Default(value bool) *BoolSchema {
	b.rules = append(b.rules, func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Value = value
		}
	})
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
