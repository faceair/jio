package main

import (
	"log"

	"github.com/faceair/jio"
)

func main() {
	data := []byte(`{
		"debug": "on",
		"list": [1, "hi"],
		"is": "yes"
	}`)
	schema := jio.Object().Keys(jio.K{
		"debug": jio.String().Valid("on", "off").Transform(func(ctx *jio.Context) error {
			if ctx.Value.(string) == "on" {
				ctx.Value = "off"
			}
			return nil
		}),
		"title": jio.String().Default("测试"),
		"list":  jio.Array().Valid("hi", jio.Number()),
		"is":    jio.Bool().Truthy("true", "yes").Required(),
	})
	err := jio.ValidateJSON(&data, schema)
	if err != nil {
		log.Printf("%v", err)
	}
	log.Printf("%s", data)
}
