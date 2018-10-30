package jio

import (
	"fmt"
	"sort"
)

type objectItem struct {
	key    string
	schema Schema
}

type K map[string]Schema

func (k K) sort() []objectItem {
	objects := make([]objectItem, 0, len(k))
	for key, schema := range k {
		objects = append(objects, objectItem{key, schema})
	}
	sort.Slice(objects, func(i, j int) bool {
		return objects[i].schema.Priority() > objects[j].schema.Priority()
	})
	return objects
}

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

func (o *ObjectSchema) SetPriority(priority int) *ObjectSchema {
	o.priority = priority
	return o
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
		fields := make([]string, len(ctx.Fields))
		copy(fields, ctx.Fields)

		for _, obj := range children.sort() {
			value, _ := ctxValue[obj.key]
			ctx.Enter(append(fields, obj.key), value)
			obj.schema.Validate(ctx)
			if ctx.skip {
				return
			}
			ctxValue[obj.key] = ctx.Value
		}
		ctx.Value = ctxValue
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
