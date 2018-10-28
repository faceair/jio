package jio

import (
	"fmt"
	"regexp"
)

func String() *StringSchema {
	return &StringSchema{
		rules: make([]func(string, string) (string, error), 0, 3),
	}
}

var _ Schema = new(StringSchema)

type StringSchema struct {
	required     *bool
	defaultValue *string
	rules        []func(string, string) (string, error)
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
	s.rules = append(s.rules, func(field string, raw string) (string, error) {
		var isValid bool
		for _, v := range values {
			if v == raw {
				isValid = true
				break
			}
		}
		if !isValid {
			return "", fmt.Errorf("field `%s` value %v is not in %v", field, raw, values)
		}
		return raw, nil
	})
	return s
}

func (s *StringSchema) Min(min int) *StringSchema {
	s.rules = append(s.rules, func(field string, raw string) (string, error) {
		if len(raw) < min {
			return "", fmt.Errorf("field `%s` value %s length less than %d", field, raw, min)
		}
		return raw, nil
	})
	return s
}

func (s *StringSchema) Max(max int) *StringSchema {
	s.rules = append(s.rules, func(field string, raw string) (string, error) {
		if len(raw) > max {
			return "", fmt.Errorf("field `%s` value %s length exceeded %d", field, raw, max)
		}
		return raw, nil
	})
	return s
}

func (s *StringSchema) Length(length int) *StringSchema {
	s.rules = append(s.rules, func(field string, raw string) (string, error) {
		if len(raw) != length {
			return "", fmt.Errorf("field `%s` value %s length not equal to %d", field, raw, length)
		}
		return raw, nil
	})
	return s
}

func (s *StringSchema) Regex(regex string) *StringSchema {
	re := regexp.MustCompile(regex)
	s.rules = append(s.rules, func(field string, raw string) (string, error) {
		if !re.MatchString(raw) {
			return "", fmt.Errorf("field `%s` value %s not match with %s", field, raw, regex)
		}
		return raw, nil
	})
	return s
}

func (s *StringSchema) Transform(f func(string, string) (string, error)) *StringSchema {
	s.rules = append(s.rules, f)
	return s
}

func (s *StringSchema) Validate(field string, raw interface{}) (interface{}, error) {
	if s.isRequired() {
		if raw == nil {
			return nil, fmt.Errorf("field `%s` is required", field)
		}
	} else {
		if raw == nil {
			if s.defaultValue != nil {
				raw = *s.defaultValue
			} else {
				return raw, nil
			}
		}
	}
	strRaw, ok := (raw).(string)
	if !ok {
		return nil, fmt.Errorf("field `%s` value %s is not string", field, raw)
	}
	var err error
	for _, rule := range s.rules {
		strRaw, err = rule(field, strRaw)
		if err != nil {
			return strRaw, err
		}
	}
	return strRaw, nil
}
