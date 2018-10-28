package jio

import (
	"fmt"
	"math"
	"strconv"
)

func Number() *NumberSchema {
	return &NumberSchema{
		rules: make([]func(string, float64) (float64, error), 0, 3),
	}
}

var _ Schema = new(NumberSchema)

type NumberSchema struct {
	required     *bool
	defaultValue *string
	rules        []func(string, float64) (float64, error)
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
	n.rules = append(n.rules, func(field string, raw float64) (float64, error) {
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
	return n
}

func (n *NumberSchema) Min(min float64) *NumberSchema {
	n.rules = append(n.rules, func(field string, raw float64) (float64, error) {
		if raw < min {
			return 0, fmt.Errorf("field `%s` value %v less than %v", field, raw, min)
		}
		return raw, nil
	})
	return n
}

func (n *NumberSchema) Max(max float64) *NumberSchema {
	n.rules = append(n.rules, func(field string, raw float64) (float64, error) {
		if raw > max {
			return 0, fmt.Errorf("field `%s` value %v exceeded %v", field, raw, max)
		}
		return raw, nil
	})
	return n
}

func (n *NumberSchema) Ceil() *NumberSchema {
	n.rules = append(n.rules, func(field string, raw float64) (float64, error) {
		return math.Ceil(raw), nil
	})
	return n
}

func (n *NumberSchema) Floor() *NumberSchema {
	n.rules = append(n.rules, func(field string, raw float64) (float64, error) {
		return math.Floor(raw), nil
	})
	return n
}

func (n *NumberSchema) Round() *NumberSchema {
	n.rules = append(n.rules, func(field string, raw float64) (float64, error) {
		return math.Floor(raw + 0.5), nil
	})
	return n
}

func (n *NumberSchema) Validate(field string, raw interface{}) (interface{}, error) {
	if n.isRequired() {
		if raw == nil {
			return nil, fmt.Errorf("field `%s` is required", field)
		}
	} else {
		if raw == nil {
			if n.defaultValue != nil {
				raw = *n.defaultValue
			} else {
				return raw, nil
			}
		}
	}
	var numRaw float64
	var err error
	switch value := (raw).(type) {
	case string:
		numRaw, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
	case float64:
		numRaw = value
	default:
		return nil, fmt.Errorf("field `%s` value %v is not number", field, raw)
	}
	for _, rule := range n.rules {
		numRaw, err = rule(field, numRaw)
		if err != nil {
			return numRaw, err
		}
	}
	return numRaw, nil
}
