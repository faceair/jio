<h1 align="center">jio</h1>

<p align="center">
    <img src="jio.jpg" width="240" height="240" border="0" alt="jio">
</p>
<p align="center">Make validation simple and efficient !</p>

<p align="center">
    <a href="https://travis-ci.org/faceair/jio"><img src="https://img.shields.io/travis/faceair/jio/master.svg" alt="Travis branch"></a>
    <a href="https://coveralls.io/github/faceair/jio?branch=master"><img src="https://coveralls.io/repos/github/faceair/jio/badge.svg?branch=master" alt="Coverage Status"></a>
    <a href="https://goreportcard.com/report/github.com/faceair/jio"><img src="https://goreportcard.com/badge/github.com/faceair/jio" alt="Go Report Card"></a>
    <a href="https://github.com/avelino/awesome-go"><img src="https://awesome.re/mentioned-badge.svg" alt="License"></a>
    <a href="https://godoc.org/github.com/faceair/jio"><img src="https://godoc.org/github.com/faceair/jio?status.svg" alt="GoDoc"></a>
    <a href="https://raw.githubusercontent.com/faceair/jio/master/LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License"></a>
</p>

[中文文档](README.zh.md)

## Why use jio?

Parameter validation in Golang is really a cursing problem. Defining tags on structs is not easy to extend rules, handwritten validation code makes logic code cumbersome, and the initial zero value of the struct field will also interfere with the validation.

jio tries validate json raw data before deserialization to avoid these problems. Defining validation rules as Schema is easy to read and easy to extend (Inspired by Hapi.js joi library). Rules within Schema can be validated in the order of registration, and context can be used to exchange data between rules, and can access other field data even within a single rule, etc.

jio provides a flexible enough way to make your validation simple and efficient!

## How to use?

### Validate json string

```go
package main

import (
    "log"

    "github.com/faceair/jio"
)

func main() {
    data := []byte(`{
        "debug": "on",
        "window": {
            "title": "Sample Widget",
            "size": [500, 500]
        }
    }`)
    _, err := jio.ValidateJSON(&data, jio.Object().Keys(jio.K{
        "debug": jio.Bool().Truthy("on").Required(),
        "window": jio.Object().Keys(jio.K{
            "title": jio.String().Min(3).Max(18),
            "size":  jio.Array().Items(jio.Number().Integer()).Length(2).Required(),
        }).Without("name", "title").Required(),
    }))
    if err != nil {
        panic(err)
    }
    log.Printf("%s", data) // {"debug":true,"window":{"size":[500,500],"title":"Sample Widget"}}
}
```

The above schema defines the following constraints:
* `debug`
    * not empty, must be a boolean value when validation end
    * allow `on` string instead of `true`
* `window`
    * not empty, object
    * not allowed for both `name` and `title`
    * The following elements exist
        * `title`
            * string, can be empty
            * length is between 3 and 18 when not empty
        * `size`
            * array, not empty
            * there are two child elements of the integer type

### Using middleware to validate request body

Take [chi](https://github.com/go-chi/chi) as an example, the other frameworks are similar.

```go
package main

import (
    "io/ioutil"
    "net/http"

    "github.com/faceair/jio"
    "github.com/go-chi/chi"
)

func main() {
    r := chi.NewRouter()
    r.Route("/people", func(r chi.Router) {
        r.With(jio.ValidateBody(jio.Object().Keys(jio.K{
            "name":  jio.String().Min(3).Max(10).Required(),
            "age":   jio.Number().Integer().Min(0).Max(100).Required(),
            "phone": jio.String().Regex(`^1[34578]\d{9}$`).Required(),
        }), jio.DefaultErrorHandler)).Post("/", func(w http.ResponseWriter, r *http.Request) {
            body, err := ioutil.ReadAll(r.Body)
            if err != nil {
                panic(err)
            }
            w.Header().Set("Content-Type", "application/json; charset=utf-8")
            w.WriteHeader(http.StatusOK)
            w.Write(body)
        })
    })
    http.ListenAndServe(":8080", r)
}
```
The second parameter of `jio.ValidateBody` is called for error handling when the validation fails.

### Validate the query parameter with middleware

```go
package main

import (
    "encoding/json"
    "net/http"

    "github.com/faceair/jio"
    "github.com/go-chi/chi"
)

func main() {
    r := chi.NewRouter()
    r.Route("/people", func(r chi.Router) {
        r.With(jio.ValidateQuery(jio.Object().Keys(jio.K{
            "keyword":  jio.String(),
            "is_adult": jio.Bool().Truthy("true", "yes").Falsy("false", "no"),
            "starts_with": jio.Number().ParseString().Integer(),
        }), jio.DefaultErrorHandler)).Get("/", func(w http.ResponseWriter, r *http.Request) {
            query := r.Context().Value(jio.ContextKeyQuery).(map[string]interface{})
            body, err := json.Marshal(query)
            if err != nil {
                panic(err)
            }
            w.Header().Set("Content-Type", "application/json; charset=utf-8")
            w.WriteHeader(http.StatusOK)
            w.Write(body)
        })
    })
    http.ListenAndServe(":8080", r)
}
```
Note that the original value of the query parameter is string, you may need to convert the value type first (for example, `jio.Number().ParseString()` or `jio.Bool().Truthy(values)`).

## API Documentation

[https://godoc.org/github.com/faceair/jio](https://godoc.org/github.com/faceair/jio)

## Advanced usage

### Workflow

Each Schema is made up of a series of rules, for example:

```go
jio.String().Min(5).Max(10).Alphanum().Lowercase()
```

In this example, String Schema has 4 rules, which are `Min(5)` `Max(10)` `Alphanum()` `Lowercase()`, will also validate in order `Min(5) ` `Max(10)` `Alphanum()` `Lowercase()`. If a rule validation fails, the Schema's validation stops and throws an error.

In order to improve the readability of the code, these three built-in rules will validate first.

* `Required()`
* `Optional()`
* `Default(value)`

For example:

```go
jio.String().Min(5).Max(10).Alphanum().Lowercase().Required()
```

The actual validation order will be `Required()` `Min(5)` `Max(10)` `Alphanum()` `Lowercase()`.

After validate all the rules, finally we check if the basic type of the data is the type of Schema. If not, the Schema will throw an error.

### Validator Context

Data transfer in the workflow depends on context, the structure is like this:

```go
Type Context struct {
    Value interface{} // Raw data, you can also reassign to change the result
}
func (ctx *Context) Ref(refPath string) (value interface{}, ok bool) { // Reference other field data
}
func (ctx *Context) Abort(err error) { // Terminate the validation and throw an error
  ...
}
func (ctx *Context) Skip() { // Skip subsequent rules
  ...
}
```

Let's try to customize a validation rule. Add a rule to use the `Transform` method:

```go
jio.String().Transform(func(ctx *jio.Context) {
    If ctx.Value != "faceair" {
        ctx.Abort(errors.New("you are not faceair"))
    }
})
```

The custom rule we added means throwing a `you are not faceair` error when the original data is not equal to `faceair`.

In fact, the built-in validation rules work in a similar way. For example, the core code of `Optional()` is:

```go
If ctx.Value == nil {
  ctx.Skip()
}
```

You can also reassign ctx.Value to change the original data. For example, the built-in `Lowercase()` converts the original string to lowercase. The core code is:

```go
ctx.Value = strings.ToLower(ctx.Value)
```

### References and Priority

In most cases, the rules only use the data of the current field, but sometimes it needs to work with other fields. For example:

```
{
    "type": "ip", // enumeration value, `ip` or `domain`
    "value": "8.8.8.8"
}
```

The validation rules of this `value` is determined by the value of `type` and can be written as

```go
jio.Object().Keys(jio.K{
        "type": jio.String().Valid("ip", "domain").SetPriority(1).Default("ip"),
        "value": jio.String().
            When("type", "ip", jio.String().Regex(`^\d+\.\d+\.\d+\.\d+$`)).
            When("type", "domain", jio.String().Regex(`^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0 -9]\.[a-zA-Z]{2,}$`)).Required(),
})
```

The `When` function can reference other field data, and if it is successful, apply the new validation rule to the current data.

In addition, you may notice that there is a `SetPriority` method in the rules of `type`. If the input data is:

```json
{
    "value": "8.8.8.8"
}
```

When the priority is not set, the validation rule of `value` may be executed first. At this time, the value of the reference `type` will be null, and the validation will fail.
Because there are validation rules that refer to each other, there may be a validation sequence requirement. When we want a field under the same Object to be validated first, we can set it to a larger priority value (default value 0).

If you want to reference data from other fields in your custom rules, you can use the `Ref` method on the context. If the referenced data is a nested object, the path to the referenced field needs to be concatenated with `.` . For example, if you want to reference `name` under `people` object then the reference path is `people.name`:

```json
{
    "type": "people",
    "people": {
        "name": "faceair"
    }
}
```

## License

[MIT](LICENSE)
