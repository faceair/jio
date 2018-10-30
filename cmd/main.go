package main

import (
	"log"

	"github.com/faceair/jio"
)

func main() {
	data := []byte(`{
		"type": null,
		"target": "www.zhihu.com."
	}`)
	schema := jio.Object().Keys(jio.K{
		"type": jio.String().Lowercase().Valid("cname", "a", "host").Default("cname").SetPriority(100),
		"target": jio.String().
			When("type", "cname", jio.String().Regex(`^.+\.$`)).
			When("type", "a", jio.String().Regex(`^\d+\.+\d+.\d+.\d+$`)).Transform(func(ctx *jio.Context) {
			value, _ := ctx.Ref("type")
			log.Printf("%v", value)
		}),
	})
	err := jio.ValidateJSON(&data, schema)
	if err != nil {
		log.Printf("%v", err)
	}
	log.Printf("%s", data)
}
