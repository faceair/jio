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
	any := Any().Transform(func(ctx *Context) {
		ctx.Abort(errors.New("2"))
	}).Transform(func(ctx *Context) {
		ctx.Abort(errors.New("3"))
	}).PrependTransform(func(ctx *Context) {
		ctx.Abort(errors.New("1"))
	}).PrependTransform(func(ctx *Context) {
		ctx.Abort(errors.New("0"))
	})
	if len(any.rules) != 4 {
		t.Error("miss function")
	}
	for i := 0; i < 4; i++ {
		ctx := NewContext(nil)
		any.rules[i](ctx)
		if ctx.err.Error() != strconv.Itoa(i) {
			t.Error("sequential error")
		}
	}
}

func TestAnySchema_Required(t *testing.T) {
	any := Any().Required()
	ctx := NewContext(nil)
	any.Validate(ctx)
	if ctx.err == nil {
		t.Error("should error when no data")
	}
}

func TestAnySchema_Optional(t *testing.T) {
	any := Any().Optional()
	ctx := NewContext(nil)
	any.Validate(ctx)
	if ctx.err != nil {
		t.Error("should no error")
	}
}

func TestAnySchema_Default(t *testing.T) {
	defaultValue := "default_value"
	any := Any().Default(defaultValue)
	ctx := NewContext(nil)
	any.Validate(ctx)
	if ctx.Value != defaultValue {
		t.Error("should set default value")
	}
}

func TestAnySchema_When(t *testing.T) {
	any := Object().Keys(K{
		"name": Any().Required(),
		"age": Any().
			When("name", "youth", Number().Min(12)).
			When("name", "adult", Number().Min(18)).
			When("name", String(), Number().Min(0)),
	})

	ctx := NewContext(map[string]interface{}{"name": "teenagers", "age": 12})
	any.Validate(ctx)
	if ctx.err != nil {
		t.Error("teenagers test failed")
	}

	ctx = NewContext(map[string]interface{}{"name": "adult", "age": 2})
	any.Validate(ctx)
	if ctx.err == nil {
		t.Error("adult test failed")
	}

	ctx = NewContext(map[string]interface{}{"name": "badcase", "age": -3})
	any.Validate(ctx)
	if ctx.err == nil {
		t.Error("badcase test failed")
	}
}

func TestAnySchema_Valid(t *testing.T) {
	any := Any().Valid("hi")

	ctx := NewContext("hi")
	any.Validate(ctx)
	if ctx.err != nil {
		t.Error("valid value test failed")
	}

	ctx = NewContext("???")
	any.Validate(ctx)
	if ctx.err == nil {
		t.Error("invalid value test failed")
	}
}

func TestAnySchema_Validate(t *testing.T) {
	any := Any()
	ctx := NewContext(nil)
	any.Validate(ctx)
	if ctx.err != nil {
		t.Error("default optional should no error")
	}
}
