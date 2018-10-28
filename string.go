package jio

import (
	"fmt"
	"regexp"
)

func String() *StringSchema {
	return &StringSchema{
		rules: make([]func(*Context) error, 0, 3),
	}
}

var _ Schema = new(StringSchema)

type StringSchema struct {
	required     *bool
	defaultValue *string
	rules        []func(*Context) error
}

func (s *StringSchema) Required() *StringSchema {
	s.required = boolPtr(true)
	return s
}

func (s *StringSchema) isRequired() bool {
	return s.required != nil && *s.required
}

func (s *StringSchema) Default(defaultValue string) *StringSchema {
	s.defaultValue = &defaultValue
	return s
}

func (s *StringSchema) Valid(values ...string) *StringSchema {
	s.rules = append(s.rules, func(ctx *Context) error {
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
	return s
}

func (s *StringSchema) Min(min int) *StringSchema {
	s.rules = append(s.rules, func(ctx *Context) error {
		if len(ctx.Value.(string)) < min {
			return fmt.Errorf("field `%s` value %s length less than %d", ctx.FieldPath(), ctx.Value, min)
		}
		return nil
	})
	return s
}

func (s *StringSchema) Max(max int) *StringSchema {
	s.rules = append(s.rules, func(ctx *Context) error {
		if len(ctx.Value.(string)) > max {
			return fmt.Errorf("field `%s` value %s length exceeded %d", ctx.FieldPath(), ctx.Value, max)
		}
		return nil
	})
	return s
}

func (s *StringSchema) Length(length int) *StringSchema {
	s.rules = append(s.rules, func(ctx *Context) error {
		if len(ctx.Value.(string)) != length {
			return fmt.Errorf("field `%s` value %s length not equal to %d", ctx.FieldPath(), ctx.Value, length)
		}
		return nil
	})
	return s
}

func (s *StringSchema) Regex(regex string) *StringSchema {
	re := regexp.MustCompile(regex)
	s.rules = append(s.rules, func(ctx *Context) error {
		if !re.MatchString(ctx.Value.(string)) {
			return fmt.Errorf("field `%s` value %s not match with %s", ctx.FieldPath(), ctx.Value, regex)
		}
		return nil
	})
	return s
}

func (s *StringSchema) Transform(f func(*Context) error) *StringSchema {
	s.rules = append(s.rules, f)
	return s
}

func (s *StringSchema) Validate(ctx *Context) (err error) {
	if s.isRequired() {
		if ctx.Value == nil {
			return fmt.Errorf("field `%s` is required", ctx.FieldPath())
		}
	} else {
		if ctx.Value == nil {
			if s.defaultValue != nil {
				ctx.Value = *s.defaultValue
			} else {
				return nil
			}
		}
	}
	if _, ok := (ctx.Value).(string); !ok {
		return fmt.Errorf("field `%s` value %s is not string", ctx.FieldPath(), ctx.Value)
	}
	for _, rule := range s.rules {
		err = rule(ctx)
		if err != nil {
			return
		}
	}
	return
}
