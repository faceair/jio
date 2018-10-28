package jio

import (
	"fmt"
)

var _ Schema = new(ArraySchema)

func Array() *ArraySchema {
	return &ArraySchema{
		rules: make([]func(*Context) error, 0, 3),
	}
}

type ArraySchema struct {
	required     *bool
	defaultValue *interface{}
	rules        []func(*Context) error
}

func (a *ArraySchema) Required() *ArraySchema {
	a.required = boolPtr(true)
	return a
}

func (a *ArraySchema) isRequired() bool {
	return a.required != nil && *a.required
}

func (a *ArraySchema) Default(value interface{}) *ArraySchema {
	a.defaultValue = &value
	return a
}

func (a *ArraySchema) Valid(values ...interface{}) *ArraySchema {
	a.rules = append(a.rules, func(ctx *Context) (err error) {
		for _, rv := range ctx.Value.([]interface{}) {
			var isValid bool
			for _, v := range values {
				if schema, ok := v.(Schema); ok {
					err = schema.Validate(&Context{ctx.Fields, rv})
					if err == nil {
						isValid = true
						break
					}
				} else {
					if v == rv {
						isValid = true
						break
					}
				}
			}
			if !isValid {
				return fmt.Errorf("field `%s` value %v is not valid type", ctx.FieldPath(), rv)
			}
		}
		return nil
	})
	return a
}

func (a *ArraySchema) Min(min int) *ArraySchema {
	a.rules = append(a.rules, func(ctx *Context) error {
		if len(ctx.Value.([]interface{})) < min {
			return fmt.Errorf("field `%s` value %s length less than %d", ctx.FieldPath(), ctx.Value, min)
		}
		return nil
	})
	return a
}

func (a *ArraySchema) Max(max int) *ArraySchema {
	a.rules = append(a.rules, func(ctx *Context) error {
		if len(ctx.Value.([]interface{})) > max {
			return fmt.Errorf("field `%s` value %s length exceeded %d", ctx.FieldPath(), ctx.Value, max)
		}
		return nil
	})
	return a
}

func (a *ArraySchema) Length(length int) *ArraySchema {
	a.rules = append(a.rules, func(ctx *Context) error {
		if len(ctx.Value.([]interface{})) != length {
			return fmt.Errorf("field `%s` value %s length not equal to %d", ctx.FieldPath(), ctx.Value, length)
		}
		return nil
	})
	return a
}

func (a *ArraySchema) Transform(f func(*Context) error) *ArraySchema {
	a.rules = append(a.rules, f)
	return a
}

func (a *ArraySchema) Validate(ctx *Context) (err error) {
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
	if _, ok := ctx.Value.([]interface{}); !ok {
		return fmt.Errorf("field `%s` value %s is not array", ctx.FieldPath(), ctx.Value)
	}
	for _, rule := range a.rules {
		err = rule(ctx)
		if err != nil {
			return
		}
	}
	return
}
