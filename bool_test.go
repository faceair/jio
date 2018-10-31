package jio

import (
	"errors"
	"strconv"
	"testing"
)

func TestBoolSchema_SetPriority(t *testing.T) {
	for _, priority := range []int{-1, 0, 100} {
		if priority != Bool().SetPriority(priority).Priority() {
			t.Error("set priority failed")
		}
	}
}

func TestBoolSchema_TransformAndPrependTransform(t *testing.T) {
	schema := Bool().Transform(func(ctx *Context) {
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
		if ctx.err.Error() != strconv.Itoa(i) {
			t.Error("sequential error")
		}
	}
}

func TestBoolSchema_Required(t *testing.T) {
	schema := Bool().Required()
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.err == nil {
		t.Error("should error when no data")
	}
}

func TestBoolSchema_Optional(t *testing.T) {
	schema := Bool().Optional()
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.err != nil {
		t.Error("should no error")
	}
}

func TestBoolSchema_Default(t *testing.T) {
	schema := Bool().Default(true)
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.Value != true {
		t.Error("should set default value")
	}
}

func TestBoolSchema_When(t *testing.T) {
	schema := Object().Keys(K{
		"bool1": Bool().Required(),
		"bool2": Bool().
			When("bool1", Bool().Equal(true), Bool().Equal(true)).
			When("bool1", false, Bool().Equal(false)),
	})

	ctx := NewContext(map[string]interface{}{"bool1": true, "bool2": true})
	schema.Validate(ctx)
	if ctx.err != nil {
		t.Errorf("bool test failed")
	}

	ctx = NewContext(map[string]interface{}{"bool1": false, "bool2": true})
	schema.Validate(ctx)
	if ctx.err == nil {
		t.Error("bool test failed")
	}

	ctx = NewContext(map[string]interface{}{"bool1": false, "bool2": false})
	schema.Validate(ctx)
	if ctx.err != nil {
		t.Error("bool test failed")
	}
}

func TestBoolSchema_Truthy(t *testing.T) {
	schema := Bool().Truthy("yes")
	ctx := NewContext("yes")
	schema.Validate(ctx)
	if ctx.Value != true {
		t.Error("should be true")
	}
}

func TestBoolSchema_Falsy(t *testing.T) {
	schema := Bool().Falsy("no")
	ctx := NewContext("no")
	schema.Validate(ctx)
	if ctx.Value != false {
		t.Error("should be false")
	}
}

func TestBoolSchema_Validate(t *testing.T) {
	schema := Bool()
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.err != nil {
		t.Error("default optional should no error")
	}
}
