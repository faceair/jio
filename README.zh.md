<h1 align="center">jio</h1>

<p align="center">
    <img src="jio.jpg" width="240" height="240" border="0" alt="jio">
</p>
<p align="center">让校验变得简单！</p>

<p align="center">
    <a href="https://raw.githubusercontent.com/faceair/jio/master/LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License"></a>
    <a href="https://travis-ci.org/faceair/jio"><img src="https://img.shields.io/travis/faceair/jio/master.svg?t=1540985641" alt="Travis branch"></a>
    <a href="https://coveralls.io/github/faceair/jio?branch=master"><img src="https://coveralls.io/repos/github/faceair/jio/badge.svg?branch=master&t=1540985641" alt="Coverage Status"></a>
    <a href="https://goreportcard.com/report/github.com/faceair/jio"><img src="https://goreportcard.com/badge/github.com/faceair/jio?t=1540985641" alt="Go Report Card"></a>
    <a href="https://godoc.org/github.com/faceair/jio"><img src="https://godoc.org/github.com/faceair/jio?status.svg" alt="GoDoc"></a>
</p>

[English DOC](README.md)

## 为什么使用 jio ？

在 Golang 中参数校验一直都是一个很头疼的问题，在 struct 上定义 tag 很难拓展规则也不够灵活，手动写校验代码会很麻烦也会让业务逻辑难以梳理，而且 struct 字段的初始零值也会对校验产生干扰。

jio 尝试在反序列化之前校验 json 原始数据来避免这些问题，将校验规则定义成 Schema 既容易阅读也很方便地拓展。Schema 内的规则可以按注册顺序校验，同时引入 context 供上下文交换数据，甚至能在单个规则内能感知其他字段数据等等。

jio 提供足够灵活的校验方式，让你的校验变得简单高效！

## 基本用法

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

## API 文档

[https://godoc.org/github.com/faceair/jio](https://godoc.org/github.com/faceair/jio)

## 高级用法

### 工作流

每一个 Schema 都是由一系列的规则组成的，例如：
```go
jio.String().Min(5).Max(10).Alphanum().Lowercase()
```
这个例子中 String Schema 共有 4 条规则，分别是  `Min(5)` `Max(10)` `Alphanum()` `Lowercase()` ，校验的顺序也会是依次校验  `Min(5)` `Max(10)` `Alphanum()` `Lowercase()`。如果中间某个规则校验失败，Schema 的校验就会停止并向外抛出错误。

同时我们为了提升代码的可读性，我们为三个内置规则设置了优先匹配的规则，分别是
* `Required()`
* `Optional()`
* `Default(value)`

即包含这三个内置规则的 Schema 将优先校验这三个规则，例如：
```go
jio.String().Min(5).Max(10).Alphanum().Lowercase().Required()
```
的实际匹配顺序将会是 `Required()` `Min(5)` `Max(10)` `Alphanum()` `Lowercase()`。

在校验执行完所有的规则后，最后我们校验数据的基本类型是否符合 Schema 的预期。如上文的例子中，如果经过所有的规则处理后最终数据的类型不是 String，那么 Schema 将会抛出错误。

### 验证上下文（Context）

工作流中的数据传递和流程控制是依靠 Context 结构完成的，略去一些内部方法和字段后的 Context 结构大概是这样的：
```go
type Context struct {
    Value    interface{}  // 需要校验的原始数据，也可以重新赋值改变结果
}
func (ctx *Context) Ref(refPath string) (value interface{}, ok bool) { // 引用其他字段数据
}
func (ctx *Context) Abort(err error) { // 终止校验并抛出错误
  ...
}
func (ctx *Context) Skip() { // 跳过这个 Schema 后续规则的校验
  ...
}
```

我们来尝试自定义一个校验规则看看 Context 是怎么使用的，添加规则可以使用 `Transform` 方法：
```go
jio.String().Transform(func(ctx *jio.Context) {
    if ctx.Value == "faceair" {
        ctx.Abort(errors.New("oh my god"))
    }
})
```
我们添加的这个自定义规则的意思是当校验数据等于 `faceair` 的时候抛出 `oh my god`的错误。

实际上内置的校验规则也是用类似的方式工作的，例如 `Optional()` 的核心代码是：
```go
if ctx.Value == nil {
  ctx.Skip()
}
```
也可以对 ctx.Value 重新赋值改变输出结果，例如内置的`Lowercase()` 是将原始字符串全部转成小写，核心代码是：
```go
ctx.Value = strings.ToLower(ctx.Value)
```

### 引用与优先级

大部分情况下的校验只用关心当前字段的数据，但也有一些时候需要跟其他字段的数据联动，例如
```
{
    "type": "ip",  // 枚举值，ip 或 domain
    "value": "8.8.8.8"
}
```
这个时候 `value` 的校验规则需要根据 `type` 的具体类型来决定，可以写成
```go
jio.Object().Keys(jio.K{
        "type": jio.String().Valid("ip", "domain").SetPriority(1).Default("ip"),
        "value": jio.String().
            When("type", "ip", jio.String().Regex(`^\d+\.\d+\.\d+\.\d+$`)).
            When("type", "domain", jio.String().Regex(`^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}$`)).Required(),
})
```
通过 `When` 函数可以实现引用其他字段数据，如果判断成功就应用新的校验规则给当前的数据。

另外你可能注意到 `type` 的校验 Schema 中添加了一个 `SetPriority` 方法。因为当校验规则存在互相引用的时候，可能会存在校验先后顺序的要求，当我们期望同一 Object 下的某个字段被优先校验的时候我们可以给它设置一个较大的优先值 (默认值优先级 0 )。
假如这里如果传入的数据为：
```json
{
    "value": "8.8.8.8"
}
```
不手动设置优先级的时候，可能 `value` 的校验规则会先执行，此时引用 `type` 的值就会是空值，无法满足我们的校验要求。

如果在自定义规则中也想引用其他字段的数据可以使用 Context 上的 `Ref` 方法。如果引用的数据是嵌套的的对象则引用字段的路径需要用 `.` 连接，例如想要引用 `people ` 下的 `name` 则引用路径为 `people.name`。
```json
{
    "type": "people",
    "people": {
        "name": "faceair"
    }
}
```
