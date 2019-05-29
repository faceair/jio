<h1 align="center">jio</h1>

<p align="center">
    <img src="jio.jpg" width="240" height="240" border="0" alt="jio">
</p>
<p align="center">让校验变得简单高效！</p>

<p align="center">
    <a href="https://travis-ci.org/faceair/jio"><img src="https://img.shields.io/travis/faceair/jio/master.svg" alt="Travis branch"></a>
    <a href="https://coveralls.io/github/faceair/jio?branch=master"><img src="https://coveralls.io/repos/github/faceair/jio/badge.svg?branch=master" alt="Coverage Status"></a>
    <a href="https://goreportcard.com/report/github.com/faceair/jio"><img src="https://goreportcard.com/badge/github.com/faceair/jio" alt="Go Report Card"></a>
    <a href="https://github.com/avelino/awesome-go"><img src="https://awesome.re/mentioned-badge.svg" alt="License"></a>
    <a href="https://godoc.org/github.com/faceair/jio"><img src="https://godoc.org/github.com/faceair/jio?status.svg" alt="GoDoc"></a>
    <a href="https://raw.githubusercontent.com/faceair/jio/master/LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License"></a>
</p>

[English DOC](README.md)

## 为什么使用 jio ？

在 Golang 中参数校验一直都是一个很头疼的问题，在 struct 上定义 tag 不好拓展规则也不够灵活，手写校验代码会让逻辑代码很啰嗦，而且 struct 字段的初始零值也会对校验产生干扰。

jio 尝试在反序列化之前校验 json 原始数据来避免这些问题，将校验规则定义成 Schema 既容易阅读也很方便地拓展 （灵感来自 Hapi.js joi 库）。Schema 内的规则可以按注册顺序校验，同时可以使用 context 在规则间交换数据，甚至能在单个规则内能访问其他字段数据等等。

jio 提供足够灵活的方式让你的校验变得简单高效！

## 怎么用？

### 校验 json 字符串

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

这个例子中定义了这些约束：
* `debug`
    * 非空，校验结束的时候必须是布尔值
    * 允许使用 `on` 字符串代替 `true`
* `window`
    * 非空，对象
    * 不允许 `name` 和 `title` 同时存在
    * 存在如下元素
        * `title`
            * 字符串，可以为空
            * 当不为空时长度在 3 到 18 之间
        * `size`
            * 数组，非空
            * 存在两个整数类型的子元素

### 使用 middleware 校验请求 body

以 [chi](https://github.com/go-chi/chi) 为例，其他的框架也是类似的。
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
校验失败时调用 `jio.ValidateBody`  的第二个参数进行错误处理。

### 使用 middleware 校验 query 参数

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
需要注意的是 query 参数的原始值都是 string，校验时可能需要先转换类型（例如 `jio.Number().ParseString()` 或 `jio.Bool().Truthy(values)`）。

## API 文档

[https://godoc.org/github.com/faceair/jio](https://godoc.org/github.com/faceair/jio)

## 高级用法

### 工作流

每一个 Schema 都是由一系列的规则组成的，例如：

```go
jio.String().Min(5).Max(10).Alphanum().Lowercase()
```

这个例子中 String Schema 共有 4 条规则，分别是  `Min(5)` `Max(10)` `Alphanum()` `Lowercase()` ，也会按顺序依次校验 `Min(5)` `Max(10)` `Alphanum()` `Lowercase()`。如果某个规则校验失败，Schema 的校验就会停止并向外抛出错误。

为了提升代码的可读性，这三个内置规则会优先匹配，分别是

* `Required()`
* `Optional()`
* `Default(value)`

例如：

```go
jio.String().Min(5).Max(10).Alphanum().Lowercase().Required()
```

的实际匹配顺序将会是 `Required()` `Min(5)` `Max(10)` `Alphanum()` `Lowercase()`。

在校验完所有的规则后，最后我们检查数据的基本类型是否是 Schema 的类型，如果不是，Schema 将会抛出错误。

### 验证上下文（Context）

工作流中的数据传递依靠 Context，结构是这样的：

```go
type Context struct {
    Value    interface{}  // 原始数据，也可以重新赋值来改变结果
}
func (ctx *Context) Ref(refPath string) (value interface{}, ok bool) { // 引用其他字段数据
}
func (ctx *Context) Abort(err error) { // 终止校验并抛出错误
  ...
}
func (ctx *Context) Skip() { // 跳过后续规则
  ...
}
```

我们来尝试自定义一个校验规则，添加规则可以使用 `Transform` 方法：

```go
jio.String().Transform(func(ctx *jio.Context) {
    if ctx.Value != "faceair" {
        ctx.Abort(errors.New("你不是 faceair"))
    }
})
```

我们添加的这个自定义规则的意思是当原始数据等于 `faceair` 的时候抛出 `你不是 faceair` 的错误。

实际上内置的校验规则也是用类似的方式工作的，例如 `Optional()` 的核心代码是：

```go
if ctx.Value == nil {
  ctx.Skip()
}
```

也可以对 ctx.Value 重新赋值改变输出结果，例如内置的 `Lowercase()` 是将原始字符串全部转成小写，核心代码是：

```go
ctx.Value = strings.ToLower(ctx.Value)
```

### 引用与优先级

大部分情况下的规则只使用当前字段的数据，但有时也需要跟其他字段配合。例如：

```
{
    "type": "ip",  // 枚举值，`ip` 或 `domain`
    "value": "8.8.8.8"
}
```

这个 `value` 的校验规则根据 `type` 的值来决定，可以写成

```go
jio.Object().Keys(jio.K{
        "type": jio.String().Valid("ip", "domain").SetPriority(1).Default("ip"),
        "value": jio.String().
            When("type", "ip", jio.String().Regex(`^\d+\.\d+\.\d+\.\d+$`)).
            When("type", "domain", jio.String().Regex(`^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}$`)).Required(),
})
```

`When` 函数可以引用其他字段数据，如果判断成功就应用新的校验规则给当前的数据。

另外，你可能注意到 `type` 的规则中有一个 `SetPriority` 方法。如果输入数据为：

```json
{
    "value": "8.8.8.8"
}
```

不设置优先级的时候，可能 `value` 的校验规则会先执行，此时引用 `type` 的值就会是空值，校验就会失败。
因为存在校验规则互相引用时，就可能会存在校验顺序的要求。当我们希望同一 Object 下的某个字段优先校验时，我们可以给它设置一个较大的优先值 (默认值优先级 0 )。

如果在自定义规则中也想引用其他字段的数据，可以使用 Context 上的 `Ref` 方法。如果引用的数据是嵌套的的对象，则引用字段的路径需要用 `.` 连接。例如，想要引用 `people` 对象下的 `name` 则引用路径为 `people.name`：

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
