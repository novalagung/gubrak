# Gubrak v2

Golang utility library with syntactic sugar. It's like lodash, but for golang.

[![Go Report Card](https://goreportcard.com/badge/github.com/novalagung/gubrak?nocache=1)](https://goreportcard.com/report/github.com/novalagung/gubrak?nocache=1)
[![Build Status](https://travis-ci.org/novalagung/gubrak.svg?branch=master)](https://travis-ci.org/novalagung/gubrak)
[![Coverage Status](https://coveralls.io/repos/github/novalagung/gubrak/badge.svg?branch=master)](https://coveralls.io/github/novalagung/gubrak?branch=master)

## Installation

The latest version of gubrak, v2, can be downloaded in three ways.

- Using `go get` from github, for `$GOPATH`-based project:

    ```go
    go get -u github.com/novalagung/gubrak
    ```

- Using `go get` from github, for **Go Mod**-based project:

    ```go
    go get -u github.com/novalagung/gubrak@v2
    ```

- Using `go get` from gopkg.in:

    ```go
    go get -u gopkg.in/novalagung/gubrak.v2
    ```

For legacy version, v1, use this:

```go
go get -u gopkg.in/novalagung/gubrak.v1
```

## Documentation

 - [Godoc](https://godoc.org/github.com/novalagung/gubrak)

## Hello World Example

```go
package main

import (
	"fmt"
	"github.com/novalagung/gubrak"
)

type Sample struct {
	EbookName      string
	DailyDownloads int
}

func main() {
	data := []Sample{
		{EbookName: "clean code", DailyDownloads: 10000},
		{EbookName: "rework", DailyDownloads: 12000},
		{EbookName: "detective comics", DailyDownloads: 11500},
	}

	result := gubrak.From(data).
		Filter(func(each Sample) bool {
			return each.DailyDownloads > 11000
		}).
		Map(func(each Sample) string {
			return each.EbookName
		}).
		Join(", ").
		Result()

	fmt.Printf("%#v \n", result) // rework, detective comics
}
```

## APIs

Below are the list of available functions on gubrak:

[Chunk](https://godoc.org/github.com/novalagung/gubrak#Chainable.Chunk) • [Compact](https://godoc.org/github.com/novalagung/gubrak#Chainable.Compact) • [ConcatMany](https://godoc.org/github.com/novalagung/gubrak#Chainable.ConcatMany) • [Concat](https://godoc.org/github.com/novalagung/gubrak#Chainable.Concat) • [CountBy](https://godoc.org/github.com/novalagung/gubrak#Chainable.CountBy) • [Count](https://godoc.org/github.com/novalagung/gubrak#Chainable.Count) • [DifferenceMany](https://godoc.org/github.com/novalagung/gubrak#Chainable.DifferenceMany) • [Difference](https://godoc.org/github.com/novalagung/gubrak#Chainable.Difference) • [Drop](https://godoc.org/github.com/novalagung/gubrak#Chainable.Drop) • [DropRight](https://godoc.org/github.com/novalagung/gubrak#Chainable.DropRight) • [Each](https://godoc.org/github.com/novalagung/gubrak#Chainable.Each) • [EachRight](https://godoc.org/github.com/novalagung/gubrak#Chainable.EachRight) • [Exclude](https://godoc.org/github.com/novalagung/gubrak#Chainable.Exclude) • [ExcludeMany](https://godoc.org/github.com/novalagung/gubrak#Chainable.ExcludeMany) • [ExcludeAt](https://godoc.org/github.com/novalagung/gubrak#Chainable.ExcludeAt) • [ExcludeAtMany](https://godoc.org/github.com/novalagung/gubrak#Chainable.ExcludeAtMany) • [Fill](https://godoc.org/github.com/novalagung/gubrak#Chainable.Fill) • [Filter](https://godoc.org/github.com/novalagung/gubrak#Chainable.Filter) • [Find](https://godoc.org/github.com/novalagung/gubrak#Chainable.Find) • [FindIndex](https://godoc.org/github.com/novalagung/gubrak#Chainable.FindIndex) • [FindLast](https://godoc.org/github.com/novalagung/gubrak#Chainable.FindLast) • [FindLastIndex](https://godoc.org/github.com/novalagung/gubrak#Chainable.FindLastIndex) • [First](https://godoc.org/github.com/novalagung/gubrak#Chainable.First) • [FromPairs](https://godoc.org/github.com/novalagung/gubrak#Chainable.FromPairs) • [GroupBy](https://godoc.org/github.com/novalagung/gubrak#Chainable.GroupBy) • [Contains](https://godoc.org/github.com/novalagung/gubrak#Chainable.Contains) • [IndexOf](https://godoc.org/github.com/novalagung/gubrak#Chainable.IndexOf) • [Initial](https://godoc.org/github.com/novalagung/gubrak#Chainable.Initial) • [Intersection](https://godoc.org/github.com/novalagung/gubrak#Chainable.Intersection) • [IntersectionMany](https://godoc.org/github.com/novalagung/gubrak#Chainable.IntersectionMany) • [Join](https://godoc.org/github.com/novalagung/gubrak#Chainable.Join) • [KeyBy](https://godoc.org/github.com/novalagung/gubrak#Chainable.KeyBy) • [Last](https://godoc.org/github.com/novalagung/gubrak#Chainable.Last) • [LastIndexOf](https://godoc.org/github.com/novalagung/gubrak#Chainable.LastIndexOf) • [Map](https://godoc.org/github.com/novalagung/gubrak#Chainable.Map) • [Nth](https://godoc.org/github.com/novalagung/gubrak#Chainable.Nth) • [OrderBy](https://godoc.org/github.com/novalagung/gubrak#Chainable.OrderBy) • [Partition](https://godoc.org/github.com/novalagung/gubrak#Chainable.Partition) • [Reduce](https://godoc.org/github.com/novalagung/gubrak#Chainable.Reduce) • [Reject](https://godoc.org/github.com/novalagung/gubrak#Chainable.Reject) • [Reverse](https://godoc.org/github.com/novalagung/gubrak#Chainable.Reverse) • [Sample](https://godoc.org/github.com/novalagung/gubrak#Chainable.Sample) • [SampleSize](https://godoc.org/github.com/novalagung/gubrak#Chainable.SampleSize) • [Shuffle](https://godoc.org/github.com/novalagung/gubrak#Chainable.Shuffle) • [Size](https://godoc.org/github.com/novalagung/gubrak#Chainable.Size) • [Tail](https://godoc.org/github.com/novalagung/gubrak#Chainable.Tail) • [Take](https://godoc.org/github.com/novalagung/gubrak#Chainable.Take) • [TakeRight](https://godoc.org/github.com/novalagung/gubrak#Chainable.TakeRight) • [Uniq](https://godoc.org/github.com/novalagung/gubrak#Chainable.Uniq) • [UnionMany](https://godoc.org/github.com/novalagung/gubrak#Chainable.UnionMany) • [IsSlice](https://godoc.org/github.com/novalagung/gubrak#IsSlice) • [IsArray](https://godoc.org/github.com/novalagung/gubrak#IsArray) • [IsSliceOrArray](https://godoc.org/github.com/novalagung/gubrak#IsSliceOrArray) • [IsBool](https://godoc.org/github.com/novalagung/gubrak#IsBool) • [IsChannel](https://godoc.org/github.com/novalagung/gubrak#IsChannel) • [IsDate](https://godoc.org/github.com/novalagung/gubrak#IsDate) • [IsString](https://godoc.org/github.com/novalagung/gubrak#IsString) • [IsEmptyString](https://godoc.org/github.com/novalagung/gubrak#IsEmptyString) • [IsFloat](https://godoc.org/github.com/novalagung/gubrak#IsFloat) • [IsFunction](https://godoc.org/github.com/novalagung/gubrak#IsFunction) • [IsInt](https://godoc.org/github.com/novalagung/gubrak#IsInt) • [IsMap](https://godoc.org/github.com/novalagung/gubrak#IsMap) • [IsNil](https://godoc.org/github.com/novalagung/gubrak#IsNil) • [IsNumeric](https://godoc.org/github.com/novalagung/gubrak#IsNumeric) • [IsPointer](https://godoc.org/github.com/novalagung/gubrak#IsPointer) • [IsStructObject](https://godoc.org/github.com/novalagung/gubrak#IsStructObject) • [IsTrue](https://godoc.org/github.com/novalagung/gubrak#IsTrue) • [IsUint](https://godoc.org/github.com/novalagung/gubrak#IsUint) • [IsZeroNumber](https://godoc.org/github.com/novalagung/gubrak#IsZeroNumber) • [IsZeroValue](https://godoc.org/github.com/novalagung/gubrak#IsZeroValue) • [RandomInt](https://godoc.org/github.com/novalagung/gubrak#RandomInt) • [RandomString](https://godoc.org/github.com/novalagung/gubrak#RandomString) • [ReplaceCaseInsensitive](https://godoc.org/github.com/novalagung/gubrak#ReplaceCaseInsensitive)

## Test

```bash
go get -u github.com/novalagung/gubrak
dep ensure
go test -cover -race -v ./... 
```

## Contribution

Fork ➜ Create branch ➜ Commit ➜ Push ➜ Pull Requests

## License

MIT License
