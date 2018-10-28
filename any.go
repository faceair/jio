package jio

import (
	"fmt"
)

var _ Schema = new(AnySchema)

func Any() *AnySchema {
	return &AnySchema{
		rules: make([]func(*Context) error, 0, 3),
	}
}

type AnySchema struct {
	required     *bool
	defaultValue *interface{}
	rules        []func(*Context) error
}

func (a *AnySchema) Required() *AnySchema {
	a.required = boolPtr(true)
	return a
}

func (a *AnySchema) isRequired() bool {
	return a.required != nil && *a.required
}

func (a *AnySchema) Default(value interface{}) *AnySchema {
	a.defaultValue = &value
	return a
}

func (a *AnySchema) Valid(values ...interface{}) *AnySchema {
	a.rules = append(a.rules, func(ctx *Context) error {
		var isValid bool
		for _, v := range values {
			if v == ctx.Value {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("field `%s` value %v is not in %v", ctx.FieldPath(), ctx.Value, values)
		}
		return nil
	})
	return a
}

func (a *AnySchema) Transform(f func(*Context) error) Schema {
	a.rules = append(a.rules, f)
	return a
}

func (a *AnySchema) Validate(ctx *Context) (err error) {
	if a.isRequired() {
		if ctx.Value == nil {
			return fmt.Errorf("field `%s` is required", ctx.FieldPath())
		}
	} else {
		if ctx.Value == nil {
			if a.defaultValue != nil {
				ctx.Value = *a.defaultValue
			} else {
				return nil
			}
		}
	}
	for _, rule := range a.rules {
		err = rule(ctx)
		if err != nil {
			return
		}
	}
	return
}
