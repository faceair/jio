package jio

import (
	"fmt"
)

type K map[string]Schema

func Object() *ObjectSchema {
	return &ObjectSchema{
		rules: make([]func(string, map[string]interface{}) (map[string]interface{}, error), 0, 3),
	}
}

var _ Schema = new(ObjectSchema)

type ObjectSchema struct {
	required     *bool
	defaultValue *map[string]interface{}
	rules        []func(string, map[string]interface{}) (map[string]interface{}, error)
}

func (o *ObjectSchema) Required() *ObjectSchema {
	o.required = boolPtr(true)
	return o
}

func (o *ObjectSchema) isRequired() bool {
	return o.required != nil && *o.required
}

func (o *ObjectSchema) Default(defaultValue map[string]interface{}) *ObjectSchema {
	o.defaultValue = &defaultValue
	return o
}

func (o *ObjectSchema) Keys(children K) *ObjectSchema {
	o.rules = append(o.rules, func(field string, jsonRaw map[string]interface{}) (map[string]interface{}, error) {
		jsonNew := make(map[string]interface{})
		for key, schema := range children {
			value, _ := jsonRaw[key]
			if len(field) != 0 {
				key = fmt.Sprintf("%s.%s", field, key)
			}
			newValue, err := schema.Validate(key, value)
			if err != nil {
				return nil, err
			}
			jsonNew[key] = newValue
		}
		return jsonNew, nil
	})
	return o
}

func (o *ObjectSchema) Transform(f func(string, map[string]interface{}) (map[string]interface{}, error)) *ObjectSchema {
	o.rules = append(o.rules, f)
	return o
}

func (o *ObjectSchema) Validate(field string, raw interface{}) (interface{}, error) {
	if o.isRequired() {
		if raw == nil {
			return nil, fmt.Errorf("field `%s` is required", field)
		}
	} else {
		if raw == nil {
			if o.defaultValue != nil {
				raw = *o.defaultValue
			} else {
				return raw, nil
			}
		}
	}
	jsonRaw, ok := (raw).(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("field `%s` value %v is not object", field, raw)
	}
	var err error
	for _, rule := range o.rules {
		jsonRaw, err = rule(field, jsonRaw)
		if err != nil {
			return jsonRaw, err
		}
	}
	return jsonRaw, nil
}
