package jio

import (
	"fmt"
	"math"
	"strconv"
)

func Number() *NumberSchema {
	return &NumberSchema{
		rules: make([]func(*Context) error, 0, 3),
	}
}

var _ Schema = new(NumberSchema)

type NumberSchema struct {
	required     *bool
	defaultValue *string
	rules        []func(*Context) error
}

func (n *NumberSchema) Required() *NumberSchema {
	n.required = boolPtr(true)
	return n
}

func (n *NumberSchema) isRequired() bool {
	return n.required != nil && *n.required
}

func (n *NumberSchema) Default(defaultValue string) *NumberSchema {
	n.defaultValue = &defaultValue
	return n
}

func (n *NumberSchema) Valid(values ...float64) *NumberSchema {
	n.rules = append(n.rules, func(ctx *Context) error {
		var isValid bool
		for _, v := range values {
			if v == ctx.Value {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("field `%s` value %v is not in %v", ctx.FieldPath(), ctx.Value, values)
		}
		return nil
	})
	return n
}

func (n *NumberSchema) Min(min float64) *NumberSchema {
	n.rules = append(n.rules, func(ctx *Context) error {
		if ctx.Value.(float64) < min {
			return fmt.Errorf("field `%s` value %v less than %v", ctx.FieldPath(), ctx.Value, min)
		}
		return nil
	})
	return n
}

func (n *NumberSchema) Max(max float64) *NumberSchema {
	n.rules = append(n.rules, func(ctx *Context) error {
		if ctx.Value.(float64) > max {
			return fmt.Errorf("field `%s` value %v exceeded %v", ctx.FieldPath(), ctx.Value, max)
		}
		return nil
	})
	return n
}

func (n *NumberSchema) Ceil() *NumberSchema {
	n.rules = append(n.rules, func(ctx *Context) error {
		ctx.Value = math.Ceil(ctx.Value.(float64))
		return nil
	})
	return n
}

func (n *NumberSchema) Floor() *NumberSchema {
	n.rules = append(n.rules, func(ctx *Context) error {
		ctx.Value = math.Floor(ctx.Value.(float64))
		return nil
	})
	return n
}

func (n *NumberSchema) Round() *NumberSchema {
	n.rules = append(n.rules, func(ctx *Context) error {
		ctx.Value = math.Floor(ctx.Value.(float64) + 0.5)
		return nil
	})
	return n
}

func (n *NumberSchema) Validate(ctx *Context) (err error) {
	if n.isRequired() {
		if ctx.Value == nil {
			return fmt.Errorf("field `%s` is required", ctx.FieldPath())
		}
	} else {
		if ctx.Value == nil {
			if n.defaultValue != nil {
				ctx.Value = *n.defaultValue
			} else {
				return nil
			}
		}
	}
	var numRaw float64
	switch value := (ctx.Value).(type) {
	case string:
		numRaw, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}

	case float64:
		numRaw = value
	default:
		return fmt.Errorf("field `%s` value %v is not number", ctx.FieldPath(), ctx.Value)
	}
	ctx.Value = numRaw
	for _, rule := range n.rules {
		err = rule(ctx)
		if err != nil {
			return
		}
	}
	return
}
