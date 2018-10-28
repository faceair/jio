package jio

import "fmt"

var _ Schema = new(ArraySchema)

func Array() *ArraySchema {
	return &ArraySchema{
		rules: make([]func(string, []interface{}) ([]interface{}, error), 0, 3),
	}
}

type ArraySchema struct {
	required     *bool
	defaultValue *interface{}
	rules        []func(string, []interface{}) ([]interface{}, error)
}

func (a *ArraySchema) Required() *ArraySchema {
	a.required = boolPtr(true)
	return a
}

func (a *ArraySchema) isRequired() bool {
	return a.required != nil && *a.required
}

func (a *ArraySchema) Default(value interface{}) *ArraySchema {
	a.defaultValue = &value
	return a
}

func (a *ArraySchema) Valid(values ...interface{}) *ArraySchema {
	a.rules = append(a.rules, func(field string, raw []interface{}) ([]interface{}, error) {
		for _, rv := range raw {
			var isValid bool
			for _, v := range values {
				if v == rv {
					isValid = true
					break
				}
			}
			if !isValid {
				return nil, fmt.Errorf("field `%s` value %v is not in %v", field, rv, values)
			}
		}
		return raw, nil
	})
	return a
}

func (a *ArraySchema) Min(min int) *ArraySchema {
	a.rules = append(a.rules, func(field string, raw []interface{}) ([]interface{}, error) {
		if len(raw) < min {
			return nil, fmt.Errorf("field `%s` value %s length less than %d", field, raw, min)
		}
		return raw, nil
	})
	return a
}

func (a *ArraySchema) Max(max int) *ArraySchema {
	a.rules = append(a.rules, func(field string, raw []interface{}) ([]interface{}, error) {
		if len(raw) > max {
			return nil, fmt.Errorf("field `%s` value %s length exceeded %d", field, raw, max)
		}
		return raw, nil
	})
	return a
}

func (a *ArraySchema) Length(length int) *ArraySchema {
	a.rules = append(a.rules, func(field string, raw []interface{}) ([]interface{}, error) {
		if len(raw) != length {
			return nil, fmt.Errorf("field `%s` value %s length not equal to %d", field, raw, length)
		}
		return raw, nil
	})
	return a
}

func (a *ArraySchema) Transform(f func(field string, raw []interface{}) ([]interface{}, error)) Schema {
	a.rules = append(a.rules, f)
	return a
}

func (a *ArraySchema) Validate(field string, raw interface{}) (interface{}, error) {
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
	arrRaw, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("field `%s` value %s is not array", field, raw)
	}
	var err error
	for _, rule := range a.rules {
		raw, err = rule(field, arrRaw)
		if err != nil {
			return raw, err
		}
	}
	return raw, nil
}
