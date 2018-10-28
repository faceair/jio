package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/faceair/jio"
)

func main() {
	data := []byte(`{
		"debug": "on",
		"list": "1,2,3",
		"is": "yes"
	}`)
	schema := jio.Object().Keys(jio.K{
		"debug": jio.String().Valid("on", "off").Transform(func(ctx *jio.Context) {
			if ctx.Value.(string) == "on" {
				ctx.Value = "off"
			}
		}),
		"title": jio.String().Default("测试"),
		"list": jio.Array().Transform(func(ctx *jio.Context) {
			values := make([]float64, 0, 3)
			strs := strings.Split(ctx.Value.(string), ",")
			for _, str := range strs {
				value, err := strconv.ParseFloat(str, 64)
				if err != nil {
					ctx.Abort(err)
					return
				}
				values = append(values, value)
			}
			ctx.Value = values
		}).Valid(1.0, 2.0, 3.0),
		"is": jio.Bool().Truthy("true", "yes").Required(),
	})
	err := jio.ValidateJSON(&data, schema)
	if err != nil {
		log.Printf("%v", err)
	}
	log.Printf("%s", data)
}
