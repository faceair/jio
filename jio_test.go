package jio

import "testing"

func TestValidateJSON(t *testing.T) {
	data := []byte(`{`)
	if err := ValidateJSON(&data, Any()); err == nil {
		t.Error("should error")
	}

	data = []byte(`{"1": 10}`)
	if err := ValidateJSON(&data, Object().Keys(K{"1": Number().Max(5)})); err == nil {
		t.Error("should error")
	}

	data = []byte(`{"1": 10}`)
	if err := ValidateJSON(&data, Object().Keys(K{"1": Any().Transform(func(ctx *Context) {
		ctx.Value = make(chan int)
	})})); err == nil {
		t.Error("should error")
	}

	data = []byte(`{"1": 10}`)
	if err := ValidateJSON(&data, Object().Keys(K{"1": Number()})); err != nil {
		t.Error("should no error")
	}
}
