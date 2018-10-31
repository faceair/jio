package jio

import (
	"errors"
	"strconv"
	"testing"
)

func TestNumberSchema_SetPriority(t *testing.T) {
	for _, priority := range []int{-1, 0, 100} {
		if priority != Number().SetPriority(priority).Priority() {
			t.Error("set priority failed")
		}
	}
}

func TestNumberSchema_TransformAndPrependTransform(t *testing.T) {
	schema := Number().Transform(func(ctx *Context) {
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

func TestNumberSchema_Required(t *testing.T) {
	schema := Number().Required()
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("should error when no data")
	}
}

func TestNumberSchema_Optional(t *testing.T) {
	schema := Number().Optional()
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("should no error")
	}
}

func TestNumberSchema_Default(t *testing.T) {
	defaultValue := 1.0
	schema := Number().Default(defaultValue)
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.Value != defaultValue {
		t.Error("should set default value")
	}
}

func TestNumberSchema_Set(t *testing.T) {
	defaultValue := 1.2
	schema := Number().Set(defaultValue)
	ctx := NewContext(2.3)
	schema.Validate(ctx)
	if ctx.Value != defaultValue {
		t.Error("should set default value")
	}
}

func TestNumberSchema_Equal(t *testing.T) {
	schema := Number().Equal(3)
	ctx := NewContext(3)
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("test equal failed")
	}

	ctx = NewContext(5)
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("test equal failed")
	}
}

func TestNumberSchema_When(t *testing.T) {
	schema := Object().Keys(K{
		"name": Any().Required(),
		"age": Number().
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

func TestNumberSchema_Check(t *testing.T) {
	schema := Number().Check(func(value float64) error {
		if value != 1.0 {
			return errors.New("not equal to 1.0")
		}
		return nil
	})
	ctx := NewContext(1.0)
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("should no error")
	}

	ctx = NewContext(2.0)
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("should error")
	}

	ctx = NewContext("???")
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("should error")
	}
}

func TestNumberSchema_Valid(t *testing.T) {
	schema := Number().Valid(1)

	ctx := NewContext(1)
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("valid value test failed")
	}

	ctx = NewContext(2)
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("invalid value test failed")
	}
}

func TestNumberSchema_Min(t *testing.T) {
	schema := Number().Min(3)
	ctx := NewContext(2)
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("test min failed")
	}

	ctx = NewContext(5)
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("test min failed")
	}
}

func TestNumberSchema_Max(t *testing.T) {
	schema := Number().Max(3)
	ctx := NewContext(2)
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("test max failed")
	}

	ctx = NewContext(5)
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("test max failed")
	}
}

func TestNumberSchema_Integer(t *testing.T) {
	schema := Number().Integer()
	ctx := NewContext(3.1)
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("test integer failed")
	}

	ctx = NewContext(5)
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("test integer failed")
	}
}

func TestNumberSchema_Convert(t *testing.T) {
	schema := Number().Convert(func(value float64) float64 {
		return value + 1
	})
	ctx := NewContext(1)
	schema.Validate(ctx)
	if ctx.Value != 2.0 {
		t.Error("test convert failed")
	}

	ctx = NewContext("??")
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("test convert failed")
	}
}

func TestNumberSchema_Ceil(t *testing.T) {
	schema := Number().Ceil()
	ctx := NewContext(1.1)
	schema.Validate(ctx)
	if ctx.Value != 2.0 {
		t.Error("test ceil failed")
	}
	ctx = NewContext(1.9)
	schema.Validate(ctx)
	if ctx.Value != 2.0 {
		t.Error("test ceil failed")
	}
}

func TestNumberSchema_Floor(t *testing.T) {
	schema := Number().Floor()
	ctx := NewContext(1.1)
	schema.Validate(ctx)
	if ctx.Value != 1.0 {
		t.Error("test floor failed")
	}
	ctx = NewContext(1.9)
	schema.Validate(ctx)
	if ctx.Value != 1.0 {
		t.Error("test floor failed")
	}
}

func TestNumberSchema_Round(t *testing.T) {
	schema := Number().Round()
	ctx := NewContext(1.1)
	schema.Validate(ctx)
	if ctx.Value != 1.0 {
		t.Error("test round failed")
	}
	ctx = NewContext(1.9)
	schema.Validate(ctx)
	if ctx.Value != 2.0 {
		t.Error("test round failed")
	}
}

func TestNumberSchema_Validate(t *testing.T) {
	schema := Number()
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("default optional should no error")
	}

	ctx = NewContext("hhh")
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("not number")
	}
}

func TestNumberSchema_ParseString(t *testing.T) {
	schema := Number().ParseString()
	ctx := NewContext("1.1")
	schema.Validate(ctx)
	if ctx.Value != 1.1 {
		t.Error("test parse string failed")
	}
	ctx = NewContext("hi1.1")
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("test parse string failed")
	}
}
