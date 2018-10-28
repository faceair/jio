package jio

import (
	"fmt"
)

type K map[string]Schema

func Object() *ObjectSchema {
	return &ObjectSchema{
		rules: make([]func(*Context), 0, 3),
	}
}

var _ Schema = new(ObjectSchema)

type ObjectSchema struct {
	rules []func(*Context)
}

func (o *ObjectSchema) Required() *ObjectSchema {
	o.rules = append(o.rules, func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Abort(fmt.Errorf("field `%s` is required", ctx.FieldPath()))
		}
	})
	return o
}

func (o *ObjectSchema) Optional() *ObjectSchema {
	o.rules = append(o.rules, func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Skip()
		}
	})
	return o
}

func (o *ObjectSchema) Default(value map[string]interface{}) *ObjectSchema {
	o.rules = append(o.rules, func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Value = value
		}
	})
	return o
}

func (o *ObjectSchema) Keys(children K) *ObjectSchema {
	o.rules = append(o.rules, func(ctx *Context) {
		ctxValue, ok := ctx.Value.(map[string]interface{})
		if !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not object", ctx.FieldPath(), ctx.Value))
			return
		}
		jsonNew := make(map[string]interface{})
		for key, schema := range children {
			value, _ := ctxValue[key]
			ctxNew := &Context{
				Fields: append(ctx.Fields, key),
				Value:  value,
			}
			schema.Validate(ctxNew)
			if ctxNew.err != nil {
				ctx.Abort(ctxNew.err)
				return
			}
			jsonNew[key] = ctxNew.Value
		}
		ctx.Value = jsonNew
	})
	return o
}

func (o *ObjectSchema) Transform(f func(*Context)) *ObjectSchema {
	o.rules = append(o.rules, f)
	return o
}

func (o *ObjectSchema) Validate(ctx *Context) {
	for _, rule := range o.rules {
		rule(ctx)
		if ctx.skip {
			return
		}
	}
	if ctx.err == nil {
		if _, ok := (ctx.Value).(map[string]interface{}); !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not object", ctx.FieldPath(), ctx.Value))
		}
	}
}
