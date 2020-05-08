package jio

import (
	"fmt"
	"sort"
	"strings"
)

type objectItem struct {
	key    string
	schema Schema
}

// K object keys schema alias
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

// Object Generates a schema object that matches object data type
func Object() *ObjectSchema {
	return &ObjectSchema{
		rules: make([]func(*Context), 0, 3),
	}
}

var _ Schema = new(ObjectSchema)

// ObjectSchema match object data type
type ObjectSchema struct {
	baseSchema

	required *bool
	rules    []func(*Context)
}

// SetPriority same as AnySchema.SetPriority
func (o *ObjectSchema) SetPriority(priority int) *ObjectSchema {
	o.priority = priority
	return o
}

// PrependTransform same as AnySchema.PrependTransform
func (o *ObjectSchema) PrependTransform(f func(*Context)) *ObjectSchema {
	o.rules = append([]func(*Context){f}, o.rules...)
	return o
}

// Transform same as AnySchema.Transform
func (o *ObjectSchema) Transform(f func(*Context)) *ObjectSchema {
	o.rules = append(o.rules, f)
	return o
}

// Required same as AnySchema.Required
func (o *ObjectSchema) Required() *ObjectSchema {
	o.required = boolPtr(true)
	return o.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Abort(fmt.Errorf("field `%s` is required", ctx.FieldPath()))
		}
	})
}

// Optional same as AnySchema.Optional
func (o *ObjectSchema) Optional() *ObjectSchema {
	o.required = boolPtr(false)
	return o.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Skip()
		}
	})
}

// Default same as AnySchema.Default
func (o *ObjectSchema) Default(value map[string]interface{}) *ObjectSchema {
	o.required = boolPtr(false)
	return o.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Value = value
		}
	})
}

// With require the presence of these keys.
func (o *ObjectSchema) With(keys ...string) *ObjectSchema {
	return o.Transform(func(ctx *Context) {
		ctxValue, ok := ctx.Value.(map[string]interface{})
		if !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not object", ctx.FieldPath(), ctx.Value))
			return
		}
		for _, key := range keys {
			_, ok := ctxValue[key]
			if !ok {
				ctx.Abort(fmt.Errorf("field `%s` not contains %v", ctx.FieldPath(), key))
				return
			}
		}
	})
}

// Without forbids the presence of these keys.
func (o *ObjectSchema) Without(keys ...string) *ObjectSchema {
	return o.Transform(func(ctx *Context) {
		ctxValue, ok := ctx.Value.(map[string]interface{})
		if !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not object", ctx.FieldPath(), ctx.Value))
			return
		}
		contains := make([]string, 0, 3)
		for _, key := range keys {
			_, ok := ctxValue[key]
			if ok {
				contains = append(contains, key)
			}
		}
		if len(contains) > 1 {
			ctx.Abort(fmt.Errorf("field `%s` contains %v", ctx.FieldPath(), strings.Join(contains, ",")))
			return
		}
	})
}

// When same as AnySchema.When
func (o *ObjectSchema) When(refPath string, condition interface{}, then Schema) *ObjectSchema {
	return o.Transform(func(ctx *Context) { o.when(ctx, refPath, condition, then) })
}

// Keys set the object keys's schema
func (o *ObjectSchema) Keys(children K) *ObjectSchema {
	return o.Transform(func(ctx *Context) {
		ctxValue, ok := ctx.Value.(map[string]interface{})
		if !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not object", ctx.FieldPath(), ctx.Value))
			return
		}
		fields := make([]string, len(ctx.fields))
		copy(fields, ctx.fields)

		defer func() {
			ctx.fields = fields
			ctx.Value = ctxValue
		}()

		for _, obj := range children.sort() {
			value, _ := ctxValue[obj.key]
			ctx.skip = false
			ctx.fields = append(fields, obj.key)
			ctx.Value = value
			obj.schema.Validate(ctx)
			if ctx.Err != nil {
				return
			}
			if !ctx.skip {
				ctxValue[obj.key] = ctx.Value
			}
		}
	})
}

// Validate same as AnySchema.Validate
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
	if ctx.Err == nil {
		if _, ok := (ctx.Value).(map[string]interface{}); !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not object", ctx.FieldPath(), ctx.Value))
		}
	}
}
