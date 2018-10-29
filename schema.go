package jio

type Schema interface {
	Validate(*Context)
}

func boolPtr(value bool) *bool {
	return &value
}

type baseSchema struct{}

func (b *baseSchema) when(ctx *Context, refPath string, condition interface{}, then Schema) {
	value, ok := ctx.Ref(refPath)
	if !ok {
		return
	}
	if conditionSchema, ok := condition.(Schema); ok {
		newCtx := NewContext(value)
		conditionSchema.Validate(newCtx)
		if newCtx.err == nil {
			then.Validate(ctx)
		}
		return
	}
	if value == condition {
		then.Validate(ctx)
	}
}
