package jio

import (
	"errors"
	"fmt"
	"math"
)

func Number() *NumberSchema {
	return &NumberSchema{
		rules: make([]func(*Context), 0, 3),
	}
}

var _ Schema = new(NumberSchema)

type NumberSchema struct {
	baseSchema

	required *bool
	rules    []func(*Context)
}

func (n *NumberSchema) SetPriority(priority int) *NumberSchema {
	n.priority = priority
	return n
}

func (n *NumberSchema) PrependTransform(f func(*Context)) *NumberSchema {
	n.rules = append([]func(*Context){f}, n.rules...)
	return n
}

func (n *NumberSchema) Transform(f func(*Context)) *NumberSchema {
	n.rules = append(n.rules, f)
	return n
}

func (n *NumberSchema) Required() *NumberSchema {
	n.required = boolPtr(true)
	return n.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Abort(fmt.Errorf("field `%s` is required", ctx.FieldPath()))
		}
	})
}

func (n *NumberSchema) Optional() *NumberSchema {
	n.required = boolPtr(false)
	return n.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Skip()
		}
	})
}

func (n *NumberSchema) Default(value float64) *NumberSchema {
	n.required = boolPtr(false)
	return n.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Value = value
		}
	})
}

func (n *NumberSchema) When(refPath string, condition interface{}, then Schema) *NumberSchema {
	return n.Transform(func(ctx *Context) { n.when(ctx, refPath, condition, then) })
}

func (n *NumberSchema) Check(f func(float64) error) *NumberSchema {
	return n.Transform(func(ctx *Context) {
		ctxValue, ok := ctx.Value.(float64)
		if !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not number", ctx.FieldPath(), ctx.Value))
			return
		}
		if err := f(ctxValue); err != nil {
			ctx.Abort(fmt.Errorf("field `%s` value %v %s", ctx.FieldPath(), ctx.Value, err.Error()))
		}
	})
}

func (n *NumberSchema) Valid(values ...float64) *NumberSchema {
	return n.Check(func(ctxValue float64) error {
		var isValid bool
		for _, v := range values {
			if v == ctxValue {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("not in %v", values)
		}
		return nil
	})
}

func (n *NumberSchema) Min(min float64) *NumberSchema {
	return n.Check(func(ctxValue float64) error {
		if ctxValue < min {
			return fmt.Errorf("less than %v", min)
		}
		return nil
	})
}

func (n *NumberSchema) Max(max float64) *NumberSchema {
	return n.Check(func(ctxValue float64) error {
		if ctxValue > max {
			return fmt.Errorf("exceeded %v", max)
		}
		return nil
	})
}

func (n *NumberSchema) Integer() *NumberSchema {
	return n.Check(func(ctxValue float64) error {
		if ctxValue != math.Trunc(ctxValue) {
			return errors.New("not integer")
		}
		return nil
	})
}

func (n *NumberSchema) Convert(f func(float64) float64) *NumberSchema {
	return n.Transform(func(ctx *Context) {
		ctxValue, ok := ctx.Value.(float64)
		if !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not number", ctx.FieldPath(), ctx.Value))
			return
		}
		ctx.Value = f(ctxValue)
	})
}

func (n *NumberSchema) Ceil() *NumberSchema {
	return n.Convert(math.Ceil)
}

func (n *NumberSchema) Floor() *NumberSchema {
	return n.Convert(math.Floor)
}

func (n *NumberSchema) Round() *NumberSchema {
	return n.Convert(func(ctxValue float64) float64 {
		return math.Floor(ctxValue + 0.5)
	})
}

func (n *NumberSchema) Validate(ctx *Context) {
	if n.required == nil {
		n.Optional()
	}
	if ctxValue, ok := ctx.Value.(int); ok {
		ctx.Value = float64(ctxValue)
	}
	for _, rule := range n.rules {
		rule(ctx)
		if ctx.skip {
			return
		}
	}
	if ctx.err == nil {
		if _, ok := (ctx.Value).(float64); !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not number", ctx.FieldPath(), ctx.Value))
		}
	}
}
