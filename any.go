package jio

import (
	"fmt"
)

var _ Schema = new(AnySchema)

// Any Generates a schema object that matches any data type
func Any() *AnySchema {
	return &AnySchema{
		rules: make([]func(*Context), 0, 3),
	}
}

// AnySchema match any data type
type AnySchema struct {
	baseSchema

	required *bool
	rules    []func(*Context)
}

// SetPriority set priority to the schema.
// A schema with a higher priority under the same object will be validate first.
func (a *AnySchema) SetPriority(priority int) *AnySchema {
	a.priority = priority
	return a
}

// PrependTransform run your transform function before othor rules.
func (a *AnySchema) PrependTransform(f func(*Context)) *AnySchema {
	a.rules = append([]func(*Context){f}, a.rules...)
	return a
}

// Transform append your transform function to rules.
func (a *AnySchema) Transform(f func(*Context)) *AnySchema {
	a.rules = append(a.rules, f)
	return a
}

// Required mark a key as required which will not allow undefined or null as value.
// All keys are optional by default.
func (a *AnySchema) Required() *AnySchema {
	a.required = boolPtr(true)
	return a.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Abort(fmt.Errorf("field `%s` is required", ctx.FieldPath()))
		}
	})
}

// Optional mark a key as optional which will allow undefined or null as values.
// When the value of the key is undefined or null, the following check rule will be skip.
// Used to annotate the schema for readability as all keys are optional by default.
func (a *AnySchema) Optional() *AnySchema {
	a.required = boolPtr(false)
	return a.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Skip()
		}
	})
}

// Default set a default value if the original value is undefined or null.
func (a *AnySchema) Default(value interface{}) *AnySchema {
	a.required = boolPtr(false)
	return a.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Value = value
		}
	})
}

// Set just set a value for the key and don't care the origin value.
func (a *AnySchema) Set(value interface{}) *AnySchema {
	return a.Transform(func(ctx *Context) {
		ctx.Value = value
	})
}

// Equal check the provided value is equal to the value of the key.
func (a *AnySchema) Equal(value interface{}) *AnySchema {
	return a.Transform(func(ctx *Context) {
		if value != ctx.Value {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not %v", ctx.FieldPath(), ctx.Value, value))
			return
		}
	})
}

// When add a conditional schema based on another key value
// The reference path support use `.` access object property, just like javascript.
// The condition can be a Schema or value.
// If condition is a schema, then this condition Schema will be used to verify the reference value.
// If condition is value, then check the condition is equal to the reference value.
// When the condition is true, the then schema will be applied to the current key value.
// Otherwise, nothing will be done.
func (a *AnySchema) When(refPath string, condition interface{}, then Schema) *AnySchema {
	return a.Transform(func(ctx *Context) { a.when(ctx, refPath, condition, then) })
}

// Valid add the provided values into the allowed whitelist and mark them as the only valid values allowed.
func (a *AnySchema) Valid(values ...interface{}) *AnySchema {
	return a.Transform(func(ctx *Context) {
		var isValid bool
		for _, v := range values {
			if v == ctx.Value {
				isValid = true
				break
			}
		}
		if !isValid {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not in %v", ctx.FieldPath(), ctx.Value, values))
			return
		}
	})
}

// Validate a value using the schema
func (a *AnySchema) Validate(ctx *Context) {
	if a.required == nil {
		a.Optional()
	}
	for _, rule := range a.rules {
		rule(ctx)
		if ctx.skip {
			return
		}
	}
}
