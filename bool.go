package jio

import "fmt"

func Bool() *BoolSchema {
	return &BoolSchema{
		rules: make([]func(*Context) error, 0, 3),
	}
}

var _ Schema = new(BoolSchema)

type BoolSchema struct {
	required     *bool
	defaultValue *bool
	rules        []func(*Context) error
}

func (b *BoolSchema) Required() *BoolSchema {
	b.required = boolPtr(true)
	return b
}

func (b *BoolSchema) isRequired() bool {
	return b.required != nil && *b.required
}

func (b *BoolSchema) Default(value bool) *BoolSchema {
	b.defaultValue = &value
	return b
}

func (b *BoolSchema) Truthy(values ...interface{}) *BoolSchema {
	b.rules = append(b.rules, func(ctx *Context) error {
		for _, v := range values {
			if v == ctx.Value {
				ctx.Value = true
				return nil
			}
		}
		return nil
	})
	return b
}

func (b *BoolSchema) Falsy(values ...interface{}) *BoolSchema {
	b.rules = append(b.rules, func(ctx *Context) error {
		for _, v := range values {
			if v == ctx.Value {
				ctx.Value = false
				return nil
			}
		}
		return nil
	})
	return b
}

func (b *BoolSchema) Validate(ctx *Context) (err error) {
	if b.isRequired() {
		if ctx.Value == nil {
			return fmt.Errorf("field `%s` is required", ctx.FieldPath())
		}
	} else {
		if ctx.Value == nil {
			if b.defaultValue != nil {
				ctx.Value = *b.defaultValue
			} else {
				return nil
			}
		}
	}
	for _, rule := range b.rules {
		err = rule(ctx)
		if err != nil {
			return
		}
	}
	if _, ok := (ctx.Value).(bool); !ok {
		return fmt.Errorf("field `%s` value %v is not boolean", ctx.FieldPath(), ctx.Value)
	}
	return
}
