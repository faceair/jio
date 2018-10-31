package jio

import (
	"errors"
	"reflect"
	"testing"
)

func TestContext_Ref(t *testing.T) {
	ctx := NewContext(map[string]interface{}{
		"1": map[string]interface{}{
			"2": "2",
		},
		"3": 3,
		"4": []int{1, 2, 3, 4},
	})
	value, ok := ctx.Ref("1")
	if !ok {
		t.Error("not found refer 1")
	}
	valueMap, ok := value.(map[string]interface{})
	if !ok {
		t.Error("type assert failed")
	}
	if len(valueMap) != 1 {
		t.Error("unknown map keys")
	}
	value, ok = ctx.Ref("1.2")
	if !ok {
		t.Error("not found refer 1.2")
	}
	if value != "2" {
		t.Error("unknown value")
	}
	value, ok = ctx.Ref("3")
	if !ok {
		t.Error("not found refer 3")
	}
	if value != 3 {
		t.Error("unknown value")
	}
	_, ok = ctx.Ref("4.1")
	if ok {
		t.Error("found refer 4")
	}
	_, ok = ctx.Ref("5")
	if ok {
		t.Error("found refer 5")
	}
}

func TestContext_FieldPath(t *testing.T) {
	ctx := NewContext(nil)
	ctx.fields = []string{"1"}
	if ctx.FieldPath() != "1" {
		t.Error("error path")
	}
	ctx.fields = []string{"1", "2"}
	if ctx.FieldPath() != "1.2" {
		t.Error("error path")
	}
}

func TestContext_Abort(t *testing.T) {
	ctx := NewContext(nil)
	ctx.Abort(errors.New("error"))
	if ctx.err == nil {
		t.Error("should have error")
	}
	if !ctx.skip {
		t.Error("should skip")
	}
}

func TestContext_Skip(t *testing.T) {
	ctx := NewContext(nil)
	ctx.Skip()
	if ctx.err != nil {
		t.Error("should no error")
	}
	if !ctx.skip {
		t.Error("should skip")
	}
}

func TestContext_SetAndGet(t *testing.T) {
	ctx := NewContext(nil)
	ctx.Set("name", "faceair")
	name, ok := ctx.Get("name")
	if !ok || name != "faceair" {
		t.Error("get failed")
	}
	name, ok = ctx.Get("age")
	if ok || name != nil {
		t.Error("should be nil")
	}
}

func toInterface(value interface{}) *interface{} {
	return &value
}

func TestContext_AssertKind(t *testing.T) {
	name := "faceair"
	ctx := NewContext(name)
	if !ctx.AssertKind(reflect.String) {
		t.Error("assert string faild")
	}
	if len(ctx.kindCache) != 1 {
		t.Error("assert string faild")
	}
	if !ctx.AssertKind(reflect.String) {
		t.Error("assert string faild")
	}
	if len(ctx.kindCache) != 1 {
		t.Error("assert string faild")
	}
}
