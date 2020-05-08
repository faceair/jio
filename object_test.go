package jio

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
)

func TestK_sort(t *testing.T) {
	schemas := K{
		"2": Any().SetPriority(2),
		"0": Any().SetPriority(0),
		"1": Any().SetPriority(1),
		"3": Any().SetPriority(3),
	}.sort()
	for i, schema := range schemas {
		if schema.key != strconv.Itoa(3-i) {
			t.Error("sort failed")
		}
	}
}

func TestObjectSchema_SetPriority(t *testing.T) {
	for _, priority := range []int{-1, 0, 100} {
		if priority != Object().SetPriority(priority).Priority() {
			t.Error("set priority failed")
		}
	}
}

func TestObjectSchema_TransformAndPrependTransform(t *testing.T) {
	schema := Object().Transform(func(ctx *Context) {
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

func TestObjectSchema_Required(t *testing.T) {
	schema := Object().Required()
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("should error when no data")
	}
}

func TestObjectSchema_Optional(t *testing.T) {
	schema := Object().Optional()

	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("should no error")
	}

	schema = Object().Keys(K{
		"hi": String(),
	})
	ctx = NewContext(map[string]interface{}{})
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("should no error")
	}
	_, ok := ctx.Value.(map[string]interface{})["hi"]
	if ok {
		t.Error("should no hi field")
	}
}

func TestObjectSchema_Default(t *testing.T) {
	defaultValue := map[string]interface{}{"1": "2"}
	schema := Object().Default(defaultValue)
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if reflect.ValueOf(ctx.Value).Len() != 1 {
		t.Error("should set default value")
	}
}

func TestObjectSchema_With(t *testing.T) {
	schema := Object().With("hi", "faceair")

	ctx := NewContext(map[string]interface{}{"hi": "11", "faceair": "111"})
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("valid value test failed")
	}

	ctx = NewContext(map[string]interface{}{"hi": "11", "othor": "111"})
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("invalid value test failed")
	}

	ctx = NewContext("hhh")
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("not map")
	}
}

func TestObjectSchema_Without(t *testing.T) {
	schema := Object().Without("hi", "faceair")

	ctx := NewContext(map[string]interface{}{"hi": "11", "faceair": "111"})
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("valid value test failed")
	}

	ctx = NewContext(map[string]interface{}{"hi": "11", "othor": "111"})
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("invalid value test failed")
	}

	ctx = NewContext("hhh")
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("not map")
	}
}

func TestObjectSchema_When(t *testing.T) {
	schema := Object().Keys(K{
		"exist": Bool().Required(),
		"object": Object().
			When("exist", true, Object().Required()).
			When("exist", false, Object().Optional()),
	})

	ctx := NewContext(map[string]interface{}{"exist": true, "object": map[string]interface{}{"1": "2"}})
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("exist test failed")
	}

	ctx = NewContext(map[string]interface{}{"exist": false, "object": nil})
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("not exist test failed")
	}

	ctx = NewContext(map[string]interface{}{"exist": "badcase", "age": -3})
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("badcase test failed")
	}
}

func TestObjectSchema_Keys(t *testing.T) {
	schema := Object().Keys(K{
		"exist": Bool().Required(),
	})

	ctx := NewContext(map[string]interface{}{"exist": true})
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("exist test failed")
	}

	ctx = NewContext("???")
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("unknown input should failed")
	}
}

func TestObjectSchema_Validate(t *testing.T) {
	schema := Object()
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("default optional should no error")
	}

	ctx = NewContext("hhh")
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("not map")
	}
}
