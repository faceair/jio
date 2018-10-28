package jio

import "encoding/json"

func ValidateJSON(dataRaw *[]byte, schema Schema) (err error) {
	jsonRaw := make(map[string]interface{})
	if err = json.Unmarshal(*dataRaw, &jsonRaw); err != nil {
		return err
	}
	ctx := NewContext(jsonRaw)
	schema.Validate(ctx)
	if ctx.err != nil {
		return ctx.err
	}
	dataNew, err := json.Marshal(ctx.Value)
	if err != nil {
		return err
	}
	*dataRaw = dataNew
	return
}
