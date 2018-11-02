package jio

import (
	"errors"
	"fmt"
	"reflect"
)

var _ Schema = new(ArraySchema)

// Array Generates a schema object that matches array data type
func Array() *ArraySchema {
	return &ArraySchema{
		rules: make([]func(*Context), 0, 3),
	}
}

// ArraySchema match array data type
type ArraySchema struct {
	baseSchema

	required *bool
	rules    []func(*Context)
}

// SetPriority same as AnySchema.SetPriority
func (a *ArraySchema) SetPriority(priority int) *ArraySchema {
	a.priority = priority
	return a
}

// PrependTransform same as AnySchema.PrependTransform
func (a *ArraySchema) PrependTransform(f func(*Context)) *ArraySchema {
	a.rules = append([]func(*Context){f}, a.rules...)
	return a
}

// Transform same as AnySchema.Transform
func (a *ArraySchema) Transform(f func(*Context)) *ArraySchema {
	a.rules = append(a.rules, f)
	return a
}

// Required same as AnySchema.Required
func (a *ArraySchema) Required() *ArraySchema {
	a.required = boolPtr(true)
	return a.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Abort(fmt.Errorf("field `%s` is required", ctx.FieldPath()))
		}
	})
}

// Optional same as AnySchema.Optional
func (a *ArraySchema) Optional() *ArraySchema {
	a.required = boolPtr(false)
	return a.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Skip()
		}
	})
}

// Default same as AnySchema.Default
func (a *ArraySchema) Default(value interface{}) *ArraySchema {
	a.required = boolPtr(false)
	return a.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Value = value
		}
	})
}

// When same as AnySchema.When
func (a *ArraySchema) When(refPath string, condition interface{}, then Schema) *ArraySchema {
	return a.Transform(func(ctx *Context) { a.when(ctx, refPath, condition, then) })
}

// Check use the provided function to validate the value of the key.
// Throws an error when the value is not a slice.
func (a *ArraySchema) Check(f func(interface{}) error) *ArraySchema {
	return a.Transform(func(ctx *Context) {
		if !ctx.AssertKind(reflect.Slice) {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not array", ctx.FieldPath(), ctx.Value))
			return
		}
		if err := f(ctx.Value); err != nil {
			ctx.Abort(fmt.Errorf("field `%s` value %v %s", ctx.FieldPath(), ctx.Value, err.Error()))
		}
	})
}

// Items check if this value can pass the validation of any schema.
func (a *ArraySchema) Items(schemas ...Schema) *ArraySchema {
	return a.Check(func(ctxValue interface{}) error {
		ctxRV := reflect.ValueOf(ctxValue)
		for i := 0; i < ctxRV.Len(); i++ {
			rv := ctxRV.Index(i).Interface()
			var isValid bool
			for _, schema := range schemas {
				ctxNew := NewContext(rv)
				schema.Validate(ctxNew)
				if ctxNew.Err == nil {
					isValid = true
					break
				}
			}
			if !isValid {
				return errors.New("not valid type")
			}
		}
		return nil
	})
}

// Min check if the length of this slice is greater than or equal to the provided length.
func (a *ArraySchema) Min(min int) *ArraySchema {
	return a.Check(func(ctxValue interface{}) error {
		if reflect.ValueOf(ctxValue).Len() < min {
			return fmt.Errorf("length less than %d", min)
		}
		return nil
	})
}

// Max check if the length of this slice is less than or equal to the provided length.
func (a *ArraySchema) Max(max int) *ArraySchema {
	return a.Check(func(ctxValue interface{}) error {
		if reflect.ValueOf(ctxValue).Len() > max {
			return fmt.Errorf("length exceeded %d", max)
		}
		return nil
	})
}

// Length check if the length of this slice is equal to the provided length.
func (a *ArraySchema) Length(length int) *ArraySchema {
	return a.Check(func(ctxValue interface{}) error {
		if reflect.ValueOf(ctxValue).Len() != length {
			return fmt.Errorf("length not equal to %d", length)
		}
		return nil
	})
}

// Validate same as AnySchema.Validate
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
	if ctx.Err == nil {
		if !ctx.AssertKind(reflect.Slice) {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not array", ctx.FieldPath(), ctx.Value))
		}
	}
}
