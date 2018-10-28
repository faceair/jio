package jio

import (
	"fmt"
)

type Schema interface {
	Validate(string, interface{}) (interface{}, error)
}

func boolPtr(value bool) *bool {
	return &value
}

var _ Schema = new(AnySchema)

func Any() *AnySchema {
	return &AnySchema{
		rules: make([]func(string, interface{}) (interface{}, error), 0, 3),
	}
}

type AnySchema struct {
	required     *bool
	defaultValue *interface{}
	rules        []func(string, interface{}) (interface{}, error)
}

func (a *AnySchema) Required() *AnySchema {
	a.required = boolPtr(true)
	return a
}

func (a *AnySchema) isRequired() bool {
	return a.required != nil && *a.required
}

func (a *AnySchema) Default(value interface{}) *AnySchema {
	a.defaultValue = &value
	return a
}

func (a *AnySchema) Valid(values ...interface{}) *AnySchema {
	a.rules = append(a.rules, func(field string, raw interface{}) (interface{}, error) {
		var isValid bool
		for _, v := range values {
			if v == raw {
				isValid = true
				break
			}
		}
		if !isValid {
			return 0, fmt.Errorf("field `%s` value %v is not in %v", field, raw, values)
		}
		return raw, nil
	})
	return a
}

func (a *AnySchema) Transform(f func(field string, raw interface{}) (interface{}, error)) Schema {
	a.rules = append(a.rules, f)
	return a
}

func (a *AnySchema) Validate(field string, raw interface{}) (interface{}, error) {
	if a.isRequired() {
		if raw == nil {
			return nil, fmt.Errorf("field `%s` is required", field)
		}
	} else {
		if raw == nil {
			if a.defaultValue != nil {
				raw = *a.defaultValue
			} else {
				return raw, nil
			}
		}
	}
	var err error
	for _, rule := range a.rules {
		raw, err = rule(field, raw)
		if err != nil {
			return raw, err
		}
	}
	return raw, nil
}
