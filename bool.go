package jio

import (
	"fmt"
)

// Bool Generates a schema object that matches bool data type
func Bool() *BoolSchema {
	return &BoolSchema{
		rules: make([]func(*Context), 0, 3),
	}
}

var _ Schema = new(BoolSchema)

// BoolSchema match bool data type
type BoolSchema struct {
	baseSchema

	required *bool
	rules    []func(*Context)
}

// SetPriority same as AnySchema.SetPriority
func (b *BoolSchema) SetPriority(priority int) *BoolSchema {
	b.priority = priority
	return b
}

// PrependTransform same as AnySchema.PrependTransform
func (b *BoolSchema) PrependTransform(f func(*Context)) *BoolSchema {
	b.rules = append([]func(*Context){f}, b.rules...)
	return b
}

// Transform same as AnySchema.Transform
func (b *BoolSchema) Transform(f func(*Context)) *BoolSchema {
	b.rules = append(b.rules, f)
	return b
}

// Required same as AnySchema.Required
func (b *BoolSchema) Required() *BoolSchema {
	b.required = boolPtr(true)
	return b.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Abort(fmt.Errorf("field `%s` is required", ctx.FieldPath()))
		}
	})
}

// Optional same as AnySchema.Optional
func (b *BoolSchema) Optional() *BoolSchema {
	b.required = boolPtr(false)
	return b.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Skip()
		}
	})
}

// Default same as AnySchema.Default
func (b *BoolSchema) Default(value bool) *BoolSchema {
	b.required = boolPtr(false)
	return b.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Value = value
		}
	})
}

// Set same as AnySchema.Set
func (b *BoolSchema) Set(value bool) *BoolSchema {
	return b.Transform(func(ctx *Context) {
		ctx.Value = value
	})
}

// Equal same as AnySchema.Equal
func (b *BoolSchema) Equal(value bool) *BoolSchema {
	return b.Transform(func(ctx *Context) {
		if value != ctx.Value {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not %v", ctx.FieldPath(), ctx.Value, value))
		}
	})
}

// When same as AnySchema.When
func (b *BoolSchema) When(refPath string, condition interface{}, then Schema) *BoolSchema {
	return b.Transform(func(ctx *Context) { b.when(ctx, refPath, condition, then) })
}

// Truthy allow for additional values to be considered valid booleans by converting them to true during validation.
func (b *BoolSchema) Truthy(values ...interface{}) *BoolSchema {
	return b.Transform(func(ctx *Context) {
		for _, v := range values {
			if v == ctx.Value {
				ctx.Value = true
			}
		}
	})
}

// Falsy allow for additional values to be considered valid booleans by converting them to false during validation.
func (b *BoolSchema) Falsy(values ...interface{}) *BoolSchema {
	return b.Transform(func(ctx *Context) {
		for _, v := range values {
			if v == ctx.Value {
				ctx.Value = false
			}
		}
	})
}

// Validate same as AnySchema.Validate
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
	if ctx.Err == nil {
		if _, ok := (ctx.Value).(bool); !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not boolean", ctx.FieldPath(), ctx.Value))
		}
	}
}
