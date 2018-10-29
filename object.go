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
	baseSchema

	required *bool
	rules    []func(*Context)
}

func (o *ObjectSchema) PrependTransform(f func(*Context)) *ObjectSchema {
	o.rules = append([]func(*Context){f}, o.rules...)
	return o
}

func (o *ObjectSchema) Transform(f func(*Context)) *ObjectSchema {
	o.rules = append(o.rules, f)
	return o
}

func (o *ObjectSchema) Required() *ObjectSchema {
	o.required = boolPtr(true)
	return o.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Abort(fmt.Errorf("field `%s` is required", ctx.FieldPath()))
		}
	})
}

func (o *ObjectSchema) Optional() *ObjectSchema {
	o.required = boolPtr(false)
	return o.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Skip()
		}
	})
}

func (o *ObjectSchema) Default(value map[string]interface{}) *ObjectSchema {
	o.required = boolPtr(false)
	return o.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Value = value
		}
	})
}

func (o *ObjectSchema) When(refPath string, condition interface{}, then Schema) *ObjectSchema {
	return o.Transform(func(ctx *Context) { o.when(ctx, refPath, condition, then) })
}

func (o *ObjectSchema) Keys(children K) *ObjectSchema {
	return o.Transform(func(ctx *Context) {
		ctxValue, ok := ctx.Value.(map[string]interface{})
		if !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not object", ctx.FieldPath(), ctx.Value))
			return
		}
		jsonNew := make(map[string]interface{})
		for key, schema := range children {
			value, _ := ctxValue[key]
			ctxNew := &Context{
				Root:   ctx.Root,
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
}

func (o *ObjectSchema) Validate(ctx *Context) {
	if o.required == nil {
		o.Optional()
	}
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
