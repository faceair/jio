package jio

import (
	"fmt"
)

type K map[string]Schema

func Object() *ObjectSchema {
	return &ObjectSchema{
		rules: make([]func(*Context) error, 0, 3),
	}
}

var _ Schema = new(ObjectSchema)

type ObjectSchema struct {
	required     *bool
	defaultValue *map[string]interface{}
	rules        []func(*Context) error
}

func (o *ObjectSchema) Required() *ObjectSchema {
	o.required = boolPtr(true)
	return o
}

func (o *ObjectSchema) isRequired() bool {
	return o.required != nil && *o.required
}

func (o *ObjectSchema) Default(defaultValue map[string]interface{}) *ObjectSchema {
	o.defaultValue = &defaultValue
	return o
}

func (o *ObjectSchema) Keys(children K) *ObjectSchema {
	o.rules = append(o.rules, func(ctx *Context) error {
		jsonNew := make(map[string]interface{})
		for key, schema := range children {
			value, _ := ctx.Value.(map[string]interface{})[key]
			ctxNew := &Context{
				Fields: append(ctx.Fields, key),
				Value:  value,
			}
			err := schema.Validate(ctxNew)
			if err != nil {
				return err
			}
			jsonNew[key] = ctxNew.Value
		}
		ctx.Value = jsonNew
		return nil
	})
	return o
}

func (o *ObjectSchema) Transform(f func(*Context) error) *ObjectSchema {
	o.rules = append(o.rules, f)
	return o
}

func (o *ObjectSchema) Validate(ctx *Context) (err error) {
	if o.isRequired() {
		if ctx.Value == nil {
			return fmt.Errorf("field `%s` is required", ctx.FieldPath())
		}
	} else {
		if ctx.Value == nil {
			if o.defaultValue != nil {
				ctx.Value = *o.defaultValue
			} else {
				return nil
			}
		}
	}
	_, ok := (ctx.Value).(map[string]interface{})
	if !ok {
		return fmt.Errorf("field `%s` value %v is not object", ctx.FieldPath(), ctx.Value)
	}
	for _, rule := range o.rules {
		err = rule(ctx)
		if err != nil {
			return
		}
	}
	return nil
}
