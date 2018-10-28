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
		"debug": jio.String().Valid("on", "off").Transform(func(_, raw string) (string, error) {
			if raw == "on" {
				return "off", nil
			}
			return raw, nil
		}),
		"title": jio.String().Default("测试"),
		"list":  jio.Array().Valid(1.0, "hi"),
		"is":    jio.Bool().Truthy("true", "yes").Required(),
	})
	err := jio.ValidateJSON(&data, schema)
	if err != nil {
		log.Printf("%v", err)
	}
	log.Printf("%s", data)
}
