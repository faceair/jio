package jio

import (
	"fmt"
	"reflect"
)

var _ Schema = new(ArraySchema)

func Array() *ArraySchema {
	return &ArraySchema{
		rules: make([]func(*Context), 0, 3),
	}
}

type ArraySchema struct {
	required *bool
	rules    []func(*Context)
}

func (a *ArraySchema) Required() *ArraySchema {
	a.required = boolPtr(true)
	a.rules = append([]func(*Context){func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Abort(fmt.Errorf("field `%s` is required", ctx.FieldPath()))
		}
	}}, a.rules...)
	return a
}

func (a *ArraySchema) Optional() *ArraySchema {
	a.required = boolPtr(false)
	a.rules = append([]func(*Context){func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Skip()
		}
	}}, a.rules...)
	return a
}

func (a *ArraySchema) Default(value []interface{}) *ArraySchema {
	a.required = boolPtr(false)
	a.rules = append([]func(*Context){func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Value = value
		}
	}}, a.rules...)
	return a
}

func (a *ArraySchema) Valid(values ...interface{}) *ArraySchema {
	a.rules = append(a.rules, func(ctx *Context) {
		if reflect.TypeOf(ctx.Value).Kind() != reflect.Slice {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not array", ctx.FieldPath(), ctx.Value))
			return
		}
		ctxRV := reflect.ValueOf(ctx.Value)
		for i := 0; i < ctxRV.Len(); i++ {
			rv := ctxRV.Index(i).Interface()
			var isValid bool
			for _, v := range values {
				if schema, ok := v.(Schema); ok {
					ctxNew := NewContext(rv)
					schema.Validate(ctxNew)
					if ctxNew.err == nil {
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
				ctx.Abort(fmt.Errorf("field `%s` value %v is not valid type", ctx.FieldPath(), rv))
			}
		}
	})
	return a
}

func (a *ArraySchema) Min(min int) *ArraySchema {
	a.rules = append(a.rules, func(ctx *Context) {
		if reflect.TypeOf(ctx.Value).Kind() != reflect.Slice {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not array", ctx.FieldPath(), ctx.Value))
			return
		}
		if reflect.ValueOf(ctx.Value).Len() < min {
			ctx.Abort(fmt.Errorf("field `%s` value %s length less than %d", ctx.FieldPath(), ctx.Value, min))
		}
	})
	return a
}

func (a *ArraySchema) Max(max int) *ArraySchema {
	a.rules = append(a.rules, func(ctx *Context) {
		if reflect.TypeOf(ctx.Value).Kind() != reflect.Slice {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not array", ctx.FieldPath(), ctx.Value))
			return
		}
		if reflect.ValueOf(ctx.Value).Len() > max {
			ctx.Abort(fmt.Errorf("field `%s` value %s length exceeded %d", ctx.FieldPath(), ctx.Value, max))
		}
	})
	return a
}

func (a *ArraySchema) Length(length int) *ArraySchema {
	a.rules = append(a.rules, func(ctx *Context) {
		if reflect.TypeOf(ctx.Value).Kind() != reflect.Slice {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not array", ctx.FieldPath(), ctx.Value))
			return
		}
		if reflect.ValueOf(ctx.Value).Len() != length {
			ctx.Abort(fmt.Errorf("field `%s` value %s length not equal to %d", ctx.FieldPath(), ctx.Value, length))
		}
	})
	return a
}

func (a *ArraySchema) Transform(f func(*Context)) *ArraySchema {
	a.rules = append(a.rules, f)
	return a
}

func (a *ArraySchema) Validate(ctx *Context) {
	if a.required == nil {
		a.Optional()
	}
	for _, rule := range a.rules {
		rule(ctx)
		if ctx.skip {
			return
		}
	}
	if ctx.err == nil {
		if reflect.TypeOf(ctx.Value).Kind() != reflect.Slice {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not array", ctx.FieldPath(), ctx.Value))
		}
	}
}
