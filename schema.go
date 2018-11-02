package jio

// Schema interface
type Schema interface {
	Priority() int
	Validate(*Context)
}

func boolPtr(value bool) *bool {
	return &value
}

type baseSchema struct {
	priority int
}

func (b *baseSchema) Priority() int {
	return b.priority
}

func (b *baseSchema) when(ctx *Context, refPath string, condition interface{}, then Schema) {
	value, ok := ctx.Ref(refPath)
	if !ok {
		return
	}
	if conditionSchema, ok := condition.(Schema); ok {
		newCtx := NewContext(value)
		conditionSchema.Validate(newCtx)
		if newCtx.Err == nil {
			then.Validate(ctx)
		}
		return
	}
	if value == condition {
		then.Validate(ctx)
	}
}
