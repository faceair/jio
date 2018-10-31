package jio

import (
	"fmt"
)

func Bool() *BoolSchema {
	return &BoolSchema{
		rules: make([]func(*Context), 0, 3),
	}
}

var _ Schema = new(BoolSchema)

type BoolSchema struct {
	baseSchema

	required *bool
	rules    []func(*Context)
}

func (b *BoolSchema) SetPriority(priority int) *BoolSchema {
	b.priority = priority
	return b
}

func (b *BoolSchema) PrependTransform(f func(*Context)) *BoolSchema {
	b.rules = append([]func(*Context){f}, b.rules...)
	return b
}

func (b *BoolSchema) Transform(f func(*Context)) *BoolSchema {
	b.rules = append(b.rules, f)
	return b
}

func (b *BoolSchema) Required() *BoolSchema {
	b.required = boolPtr(true)
	return b.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Abort(fmt.Errorf("field `%s` is required", ctx.FieldPath()))
		}
	})
}

func (b *BoolSchema) Optional() *BoolSchema {
	b.required = boolPtr(false)
	return b.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Skip()
		}
	})
}

func (b *BoolSchema) Default(value bool) *BoolSchema {
	b.required = boolPtr(false)
	return b.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Value = value
		}
	})
}

func (b *BoolSchema) Set(value bool) *BoolSchema {
	return b.Transform(func(ctx *Context) {
		ctx.Value = value
	})
}

func (b *BoolSchema) Equal(value bool) *BoolSchema {
	return b.Transform(func(ctx *Context) {
		if value != ctx.Value {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not %v", ctx.FieldPath(), ctx.Value, value))
			return
		}
	})
}

func (b *BoolSchema) When(refPath string, condition interface{}, then Schema) *BoolSchema {
	return b.Transform(func(ctx *Context) { b.when(ctx, refPath, condition, then) })
}

func (b *BoolSchema) Truthy(values ...interface{}) *BoolSchema {
	return b.Transform(func(ctx *Context) {
		for _, v := range values {
			if v == ctx.Value {
				ctx.Value = true
			}
		}
	})
}

func (b *BoolSchema) Falsy(values ...interface{}) *BoolSchema {
	return b.Transform(func(ctx *Context) {
		for _, v := range values {
			if v == ctx.Value {
				ctx.Value = false
			}
		}
	})
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
