package jio

import (
	"errors"
	"strconv"
	"testing"
)

func TestStringSchema_SetPriority(t *testing.T) {
	for _, priority := range []int{-1, 0, 100} {
		if priority != String().SetPriority(priority).Priority() {
			t.Error("set priority failed")
		}
	}
}

func TestStringSchema_TransformPrependTransform(t *testing.T) {
	schema := String().Transform(func(ctx *Context) {
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

func TestStringSchema_Required(t *testing.T) {
	schema := String().Required()
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("should error when no data")
	}
}

func TestStringSchema_Optional(t *testing.T) {
	schema := String().Optional()
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("should no error")
	}
}

func TestStringSchema_Default(t *testing.T) {
	defaultValue := "hi"
	schema := String().Default(defaultValue)
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.Value != defaultValue {
		t.Error("should set default value")
	}
}

func TestStringSchema_Set(t *testing.T) {
	defaultValue := "hi"
	schema := String().Set(defaultValue)
	ctx := NewContext("???")
	schema.Validate(ctx)
	if ctx.Value != defaultValue {
		t.Error("should set default value")
	}
}

func TestStringSchema_Equal(t *testing.T) {
	schema := String().Equal("faceair")
	ctx := NewContext("faceair")
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("test equal failed")
	}

	ctx = NewContext("unknown")
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("test equal failed")
	}
}

func TestStringSchema_When(t *testing.T) {
	schema := Object().Keys(K{
		"name": String().
			When("age", Number().Min(18), String().Set("adult")).
			When("age", Number().Max(17), String().Set("teenagers")).Required(),
		"age": Number().Required().SetPriority(1),
	})

	ctx := NewContext(map[string]interface{}{"age": 12, "name": "unknown"})
	schema.Validate(ctx)
	if ctx.Value.(map[string]interface{})["name"] != "teenagers" {
		t.Error("teenagers test failed")
	}

	ctx = NewContext(map[string]interface{}{"age": 20, "name": "unknown"})
	schema.Validate(ctx)
	if ctx.Value.(map[string]interface{})["name"] != "adult" {
		t.Error("adult test failed")
	}
}

func TestStringSchema_Check(t *testing.T) {
	schema := String().Check(func(value string) error {
		if value != "faceair" {
			return errors.New("not equal to faceair")
		}
		return nil
	})
	ctx := NewContext("faceair")
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("should no error")
	}

	ctx = NewContext("unknown")
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("should error")
	}

	ctx = NewContext(121213)
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("should error")
	}
}

func TestStringSchema_Valid(t *testing.T) {
	schema := String().Valid("faceair")

	ctx := NewContext("faceair")
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

func TestStringSchema_Min(t *testing.T) {
	schema := String().Min(3)
	ctx := NewContext("1234")
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("test min failed")
	}

	ctx = NewContext("1")
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("test min failed")
	}
}

func TestStringSchema_Max(t *testing.T) {
	schema := String().Max(3)
	ctx := NewContext("1")
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("test max failed")
	}

	ctx = NewContext("23333")
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("test max failed")
	}
}

func TestStringSchema_Length(t *testing.T) {
	schema := String().Length(3)
	ctx := NewContext("123")
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("test max failed")
	}

	ctx = NewContext("23333")
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("test max failed")
	}
}

func TestStringSchema_Regex(t *testing.T) {
	schema := String().Regex(`^.+\.$`)
	ctx := NewContext("google.com.")
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("test regex failed")
	}

	ctx = NewContext("google.com")
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("test regex failed")
	}
}

func TestStringSchema_Alphanum(t *testing.T) {
	schema := String().Alphanum()
	ctx := NewContext("google")
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("test alphanum failed")
	}

	ctx = NewContext("google.com")
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("test alphanum failed")
	}
}

func TestStringSchema_Token(t *testing.T) {
	schema := String().Token()
	ctx := NewContext("xsoi2n1ks_")
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("test token failed")
	}

	ctx = NewContext("hi faceair")
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("test token failed")
	}
}

func TestStringSchema_Convert(t *testing.T) {
	schema := String().Convert(func(value string) string {
		return value + "111"
	})
	ctx := NewContext("h")
	schema.Validate(ctx)
	if ctx.Value != "h111" {
		t.Error("test convert failed")
	}

	ctx = NewContext(1213213)
	schema.Validate(ctx)
	if ctx.Err == nil {
		t.Error("test convert failed")
	}
}

func TestStringSchema_Lowercase(t *testing.T) {
	schema := String().Lowercase()
	ctx := NewContext("fACeAIr")
	schema.Validate(ctx)
	if ctx.Value != "faceair" {
		t.Error("test lowercase failed")
	}
}

func TestStringSchema_Uppercase(t *testing.T) {
	schema := String().Uppercase()
	ctx := NewContext("fACeAIr")
	schema.Validate(ctx)
	if ctx.Value != "FACEAIR" {
		t.Error("test uppercase failed")
	}
}

func TestStringSchema_Trim(t *testing.T) {
	schema := String().Trim()
	ctx := NewContext("   faceair 		")
	schema.Validate(ctx)
	if ctx.Value != "faceair" {
		t.Error("test trim failed")
	}
}

func TestStringSchema_Validate(t *testing.T) {
	schema := String()
	ctx := NewContext(nil)
	schema.Validate(ctx)
	if ctx.Err != nil {
		t.Error("default optional should no error")
	}
}
