# GUBRAK

Golang utility library with syntactic sugar. It's like lodash, but for golang.

[![Go Report Card](https://goreportcard.com/badge/github.com/novalagung/gubrak?nocache=1)](https://goreportcard.com/report/github.com/novalagung/gubrak?nocache=1)
[![Build Status](https://travis-ci.org/novalagung/gubrak.svg?branch=master)](https://travis-ci.org/novalagung/gubrak)
[![cover.run](https://cover.run/go/github.com/novalagung/gubrak.svg?style=flat&tag=golang-1.10)](https://cover.run/go?tag=golang-1.10&repo=github.com%2Fnovalagung%2Fgubrak)
<!-- [![Coverage Status](https://coveralls.io/repos/github/novalagung/gubrak/badge.svg?branch=master)](https://coveralls.io/github/novalagung/gubrak?branch=master) -->

Gubrak is yet another utility library for Golang, inspired from lodash. Currently we have around 73 reusable functions available, we'll definitely adding more!

## Installation

```go
go get -u github.com/novalagung/gubrak
```

## Documentation

 - [API Documentation](https://gubrak.github.io/)
 - [Godoc](https://godoc.org/github.com/novalagung/gubrak)

## Hello World Example

```go
package main

import (
    "github.com/novalagung/gubrak"
    "fmt"
)

type Sample struct {
    EbookName      string
    DailyDownloads int
}

func main() {
    data := []Sample{
        { EbookName: "clean code", DailyDownloads: 10000 },
        { EbookName: "rework", DailyDownloads: 12000 },
        { EbookName: "detective comics", DailyDownloads: 11500 },
    }

    result, err := gubrak.Filter(data, func(each Sample) bool {
        return each.DailyDownloads > 11000
    })

    if err != nil {
        fmt.Println("Error!", err.Error)
        return
    }

    fmt.Printf("%#v \n", result.([]Sample))

    /*
    []Sample{
        { EbookName: "rework", DailyDownloads: 12000 },
        { EbookName: "detective comics", DailyDownloads: 11500 },
    }
    */
}
```

## APIs

Below are the list of available functions on gubrak:

 gubrak.Chunk() gubrak.Compact() gubrak.Concat() gubrak.Count() gubrak.Difference() gubrak.Drop() gubrak.DropRight() gubrak.Each() gubrak.EachRight() gubrak.Fill() gubrak.Filter() gubrak.Find() gubrak.FindIndex() gubrak.FindLast() gubrak.FindLastIndex() gubrak.First() gubrak.ForEach() gubrak.ForEachRight() gubrak.FromPairs() gubrak.GroupBy() gubrak.Head() gubrak.Includes() gubrak.IndexOf() gubrak.Initial() gubrak.Intersection() gubrak.IsArray() gubrak.IsBool() gubrak.IsChannel() gubrak.IsDate() gubrak.IsEmpty() gubrak.IsEmptyString() gubrak.IsFloat() gubrak.IsFunction() gubrak.IsInt() gubrak.IsMap() gubrak.IsNil() gubrak.IsNumeric() gubrak.IsPointer() gubrak.IsSlice() gubrak.IsString() gubrak.IsStructObject() gubrak.IsTrue() gubrak.IsUint() gubrak.IsZeroNumber() gubrak.Join() gubrak.KeyBy() gubrak.Last() gubrak.LastIndexOf() gubrak.Map() gubrak.Now() gubrak.Nth() gubrak.OrderBy() gubrak.Partition() gubrak.Pull() gubrak.PullAll() gubrak.PullAt() gubrak.RandomInt() gubrak.RandomString() gubrak.Reduce() gubrak.Reject() gubrak.Remove() gubrak.Reverse() gubrak.Sample() gubrak.SampleSize() gubrak.Shuffle() gubrak.Size() gubrak.SortBy() gubrak.Tail() gubrak.Take() gubrak.TakeRight() gubrak.Union() gubrak.Uniq() gubrak.Without()

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
