package jio

import (
	"fmt"
	"regexp"
	"strings"
)

func String() *StringSchema {
	return &StringSchema{
		rules: make([]func(*Context), 0, 3),
	}
}

var _ Schema = new(StringSchema)

type StringSchema struct {
	baseSchema

	required *bool
	rules    []func(*Context)
}

func (s *StringSchema) SetPriority(priority int) *StringSchema {
	s.priority = priority
	return s
}

func (s *StringSchema) PrependTransform(f func(*Context)) *StringSchema {
	s.rules = append([]func(*Context){f}, s.rules...)
	return s
}

func (s *StringSchema) Transform(f func(*Context)) *StringSchema {
	s.rules = append(s.rules, f)
	return s
}

func (s *StringSchema) Required() *StringSchema {
	s.required = boolPtr(true)
	return s.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Abort(fmt.Errorf("field `%s` is required", ctx.FieldPath()))
		}
	})
}

func (s *StringSchema) Optional() *StringSchema {
	s.required = boolPtr(false)
	return s.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Skip()
		}
	})
}

func (s *StringSchema) Default(value string) *StringSchema {
	s.required = boolPtr(false)
	return s.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Value = value
		}
	})
}

func (s *StringSchema) Set(value string) *StringSchema {
	return s.Transform(func(ctx *Context) {
		ctx.Value = value
	})
}

func (s *StringSchema) Equal(value string) *StringSchema {
	return s.Check(func(ctxValue string) error {
		if value != ctxValue {
			return fmt.Errorf("is not %v", value)
		}
		return nil
	})
}

func (s *StringSchema) When(refPath string, condition interface{}, then Schema) *StringSchema {
	return s.Transform(func(ctx *Context) { s.when(ctx, refPath, condition, then) })
}

func (s *StringSchema) Check(f func(string) error) *StringSchema {
	return s.Transform(func(ctx *Context) {
		ctxValue, ok := ctx.Value.(string)
		if !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not string", ctx.FieldPath(), ctx.Value))
			return
		}
		if err := f(ctxValue); err != nil {
			ctx.Abort(fmt.Errorf("field `%s` value %v %s", ctx.FieldPath(), ctx.Value, err.Error()))
		}
	})
}

func (s *StringSchema) Valid(values ...string) *StringSchema {
	return s.Check(func(ctxValue string) error {
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

func (s *StringSchema) Min(min int) *StringSchema {
	return s.Check(func(ctxValue string) error {
		if len(ctxValue) < min {
			return fmt.Errorf("length less than %d", min)
		}
		return nil
	})
}

func (s *StringSchema) Max(max int) *StringSchema {
	return s.Check(func(ctxValue string) error {
		if len(ctxValue) > max {
			return fmt.Errorf("length exceeded %d", max)
		}
		return nil
	})
}

func (s *StringSchema) Length(length int) *StringSchema {
	return s.Check(func(ctxValue string) error {
		if len(ctxValue) != length {
			return fmt.Errorf("length not equal to %d", length)
		}
		return nil
	})
}

func (s *StringSchema) Regex(regex string) *StringSchema {
	re := regexp.MustCompile(regex)
	return s.Check(func(ctxValue string) error {
		if !re.MatchString(ctxValue) {
			return fmt.Errorf("not match with %s", regex)
		}
		return nil
	})
}

func (s *StringSchema) Alphanum() *StringSchema {
	return s.Regex(`^[a-zA-Z0-9]+$`)
}

func (s *StringSchema) Token() *StringSchema {
	return s.Regex(`^\w+$`)
}

func (s *StringSchema) Email() *StringSchema {
	return s.Regex("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
}

func (s *StringSchema) Convert(f func(string) string) *StringSchema {
	return s.Transform(func(ctx *Context) {
		ctxValue, ok := ctx.Value.(string)
		if !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not string", ctx.FieldPath(), ctx.Value))
			return
		}
		ctx.Value = f(ctxValue)
	})
}

func (s *StringSchema) Lowercase() *StringSchema {
	return s.Convert(strings.ToLower)

}
func (s *StringSchema) Uppercase() *StringSchema {
	return s.Convert(strings.ToUpper)
}

func (s *StringSchema) Trim() *StringSchema {
	return s.Convert(strings.TrimSpace)
}

func (s *StringSchema) Validate(ctx *Context) {
	if s.required == nil {
		s.Optional()
	}
	for _, rule := range s.rules {
		rule(ctx)
		if ctx.skip {
			return
		}
	}
	if ctx.Err == nil {
		if _, ok := (ctx.Value).(string); !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not string", ctx.FieldPath(), ctx.Value))
		}
	}
}
