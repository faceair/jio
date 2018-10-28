package jio

import (
	"fmt"
	"regexp"
)

func String() *StringSchema {
	return &StringSchema{
		rules: make([]func(*Context), 0, 3),
	}
}

var _ Schema = new(StringSchema)

type StringSchema struct {
	rules []func(*Context)
}

func (s *StringSchema) Required() *StringSchema {
	s.rules = append(s.rules, func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Abort(fmt.Errorf("field `%s` is required", ctx.FieldPath()))
		}
	})
	return s
}

func (s *StringSchema) Optional() *StringSchema {
	s.rules = append(s.rules, func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Skip()
		}
	})
	return s
}

func (s *StringSchema) Default(value string) *StringSchema {
	s.rules = append(s.rules, func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Value = value
		}
	})
	return s
}

func (s *StringSchema) Valid(values ...string) *StringSchema {
	s.rules = append(s.rules, func(ctx *Context) {
		var isValid bool
		for _, v := range values {
			if v == ctx.Value {
				isValid = true
				break
			}
		}
		if !isValid {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not in %v", ctx.FieldPath(), ctx.Value, values))
		}
	})
	return s
}

func (s *StringSchema) Min(min int) *StringSchema {
	s.rules = append(s.rules, func(ctx *Context) {
		ctxValue, ok := ctx.Value.(string)
		if !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not string", ctx.FieldPath(), ctx.Value))
			return
		}
		if len(ctxValue) < min {
			ctx.Abort(fmt.Errorf("field `%s` value %s length less than %d", ctx.FieldPath(), ctx.Value, min))
		}
	})
	return s
}

func (s *StringSchema) Max(max int) *StringSchema {
	s.rules = append(s.rules, func(ctx *Context) {
		ctxValue, ok := ctx.Value.(string)
		if !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not string", ctx.FieldPath(), ctx.Value))
			return
		}
		if len(ctxValue) > max {
			ctx.Abort(fmt.Errorf("field `%s` value %s length exceeded %d", ctx.FieldPath(), ctx.Value, max))
		}
	})
	return s
}

func (s *StringSchema) Length(length int) *StringSchema {
	s.rules = append(s.rules, func(ctx *Context) {
		ctxValue, ok := ctx.Value.(string)
		if !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not string", ctx.FieldPath(), ctx.Value))
			return
		}
		if len(ctxValue) != length {
			ctx.Abort(fmt.Errorf("field `%s` value %s length not equal to %d", ctx.FieldPath(), ctx.Value, length))
		}
	})
	return s
}

func (s *StringSchema) Regex(regex string) *StringSchema {
	re := regexp.MustCompile(regex)
	s.rules = append(s.rules, func(ctx *Context) {
		ctxValue, ok := ctx.Value.(string)
		if !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not string", ctx.FieldPath(), ctx.Value))
			return
		}
		if !re.MatchString(ctxValue) {
			ctx.Abort(fmt.Errorf("field `%s` value %s not match with %s", ctx.FieldPath(), ctx.Value, regex))
		}
	})
	return s
}

func (s *StringSchema) Transform(f func(*Context)) *StringSchema {
	s.rules = append(s.rules, f)
	return s
}

func (s *StringSchema) Validate(ctx *Context) {
	for _, rule := range s.rules {
		rule(ctx)
		if ctx.skip {
			return
		}
	}
	if ctx.err == nil {
		if _, ok := (ctx.Value).(string); !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not string", ctx.FieldPath(), ctx.Value))
		}
	}
}
