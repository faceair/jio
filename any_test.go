package jio

import (
	"errors"
	"strconv"
	"testing"
)

func TestAnySchema_SetPriority(t *testing.T) {
	for _, priority := range []int{-1, 0, 100} {
		if priority != Any().SetPriority(priority).Priority() {
			t.Error("set priority failed")
		}
	}
}

func TestAnySchema_TransformAndPrependTransform(t *testing.T) {
	schema := Any().Transform(func(ctx *Context) {
		ctx.Abort(errors.New("2"))
	}).Transform(func(ctx *Context) {
		ctx.Abort(errors.New("3"))
	}).PrependTransform(func(ctx *Context) {
		ctx.Abort(errors.New("1"))
	}).PrependTransform(func(ctx *Context) {
		ctx.Abort(errors.New("0"))
	})
	if len(schema.rules) != 4 {
		t.Error("miss function")
	}
	for i := 0; i < 4; i++ {
		ctx := NewContext(nil)
		schema.rules[i](ctx)
		if ctx.Err.Error() != strconv.Itoa(i) {
			t.Error("sequential error")
		}
	}
}

func TestAnySchema_Required(t *testing.T) {
	schema := Any().Required()
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("should error when no data")
	}
}

func TestAnySchema_Optional(t *testing.T) {
	schema := Any().Optional()
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("should no error")
	}
}

func TestAnySchema_Default(t *testing.T) {
	defaultValue := "default_value"
	schema := Any().Default(defaultValue)
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.Value != defaultValue {
		t.Error("should set default value")
	}
}

func TestAnySchema_Set(t *testing.T) {
	defaultValue := "default_value"
	schema := Any().Set(defaultValue)
	ctx := NewContext("othor_value")
	schema.Validate(ctx)
	if ctx.Value != defaultValue {
		t.Error("should set default value")
	}
}

func TestAnySchema_Equal(t *testing.T) {
	schema := Any().Equal("hi")

	ctx := NewContext("hi")
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("equal value test failed")
	}

	ctx = NewContext("???")
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("equal value test failed")
	}
}

func TestAnySchema_When(t *testing.T) {
	schema := Object().Keys(K{
		"name": Any().Required(),
		"age": Any().
			When("name", "youth", Number().Min(12)).
			When("name", "adult", Number().Min(18)).
			When("name", String(), Number().Min(0)).
			When("???", String(), Number().Min(-10)),
	})

	ctx := NewContext(map[string]interface{}{"name": "teenagers", "age": 12})
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("teenagers test failed")
	}

	ctx = NewContext(map[string]interface{}{"name": "adult", "age": 2})
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("adult test failed")
	}

	ctx = NewContext(map[string]interface{}{"name": "badcase", "age": -3})
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("badcase test failed")
	}
}

func TestAnySchema_Valid(t *testing.T) {
	schema := Any().Valid("hi")

	ctx := NewContext("hi")
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("valid value test failed")
	}

	ctx = NewContext("???")
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("invalid value test failed")
	}
}

func TestAnySchema_Validate(t *testing.T) {
	schema := Any()
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("default optional should no error")
	}
}
