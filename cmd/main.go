package main

import (
	"log"

	"github.com/faceair/jio"
)

func main() {
	data := []byte(`{
		"type": "cname",
		"target": "www.zhihu.com."
	}`)
	schema := jio.Object().Keys(jio.K{
		"type": jio.String().Lowercase().Valid("cname", "a", "host").Required(),
		"target": jio.String().
			When("type", "cname", jio.String().Regex(`^.+\.$`)).
			When("type", "a", jio.String().Regex(`^\d+\.+\d+.\d+.\d+$`)),
	})
	err := jio.ValidateJSON(&data, schema)
	if err != nil {
		log.Printf("%v", err)
	}
	log.Printf("%s", data)
}
