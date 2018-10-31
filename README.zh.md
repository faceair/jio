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

jio 尝试在序列化之前校验 json 原始数据来避免这些问题，同时使用定义式的组织结构既容易阅读也能很方便的拓展规则，如果你使用过 [joi](https://github.com/hapijs/joi) 你就会感受到这种校验方式的魅力。

jio 实现了 schema 内按规则的注册顺序校验，同时引入了 context 来供校验规则的上下文交换数据，甚至在单个规则内能感知其他 schema 字段数据等等。jio 功能非常强大和灵活，请放心，jio 一定能帮你更好地校验你的数据。

## 例子

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
- `debug`
  - 非空，校验结束的时候必须是布尔值
  - 允许使用 `on` 字符串代替 `true`
- `window`
  - 非空，对象
  - 不允许 `name` 和 `title` 同时存在
  - 存在如下元素
    - `title`
      - 字符串，可以为空
      - 当不为空时长度在 3 到 18 之间
    - `size`
      - 数组，非空
      - 存在两个整数类型的子元素
