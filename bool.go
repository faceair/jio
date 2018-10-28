package jio

import "fmt"

func Bool() *BoolSchema {
	return &BoolSchema{
		rules: make([]func(string, interface{}) (interface{}, error), 0, 3),
	}
}

var _ Schema = new(BoolSchema)

type BoolSchema struct {
	required     *bool
	defaultValue *bool
	rules        []func(string, interface{}) (interface{}, error)
}

func (b *BoolSchema) Required() *BoolSchema {
	b.required = boolPtr(true)
	return b
}

func (b *BoolSchema) isRequired() bool {
	return b.required != nil && *b.required
}

func (b *BoolSchema) Default(value bool) *BoolSchema {
	b.defaultValue = &value
	return b
}

func (b *BoolSchema) Truthy(values ...interface{}) *BoolSchema {
	b.rules = append(b.rules, func(field string, raw interface{}) (interface{}, error) {
		for _, v := range values {
			if v == raw {
				return true, nil
			}
		}
		return raw, nil
	})
	return b
}

func (b *BoolSchema) Falsy(values ...interface{}) *BoolSchema {
	b.rules = append(b.rules, func(field string, raw interface{}) (interface{}, error) {
		for _, v := range values {
			if v == raw {
				return false, nil
			}
		}
		return raw, nil
	})
	return b
}

func (b *BoolSchema) Validate(field string, raw interface{}) (interface{}, error) {
	if b.isRequired() {
		if raw == nil {
			return nil, fmt.Errorf("field `%s` is required", field)
		}
	} else {
		if raw == nil {
			if b.defaultValue != nil {
				raw = *b.defaultValue
			} else {
				return raw, nil
			}
		}
	}
	var err error
	for _, rule := range b.rules {
		raw, err = rule(field, raw)
		if err != nil {
			return raw, err
		}
	}
	boolRaw, ok := (raw).(bool)
	if !ok {
		return nil, fmt.Errorf("field `%s` value %v is not boolean", field, raw)
	}
	return boolRaw, nil
}
