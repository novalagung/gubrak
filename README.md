# GUBRAK

Golang utility library with syntactic sugar. It's like lodash, but for golang.

[![Go Report Card](https://goreportcard.com/badge/github.com/novalagung/gubrak?nocache=1)](https://goreportcard.com/report/github.com/novalagung/gubrak?nocache=1)
[![Build Status](https://travis-ci.org/novalagung/gubrak.svg?branch=master)](https://travis-ci.org/novalagung/gubrak)
[![Coverage Status](https://coveralls.io/repos/github/novalagung/gubrak/badge.svg?branch=master)](https://coveralls.io/github/novalagung/gubrak?branch=master)

Gubrak is yet another utility library for Golang, inspired from lodash. Currently we have around 73 reusable functions available, we'll definitely adding more!

## Installation

```go
go get -u github.com/novalagung/gubrak
```

## Documentation

 - [API Documentation](https://gubrak.github.io/) (always updated)
 - [Godoc](https://godoc.org/github.com/novalagung/gubrak) (update delayed)

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

[Chunk](https://godoc.org/github.com/novalagung/gubrak#Chunk) • [Compact](https://godoc.org/github.com/novalagung/gubrak#Compact) • [Concat](https://godoc.org/github.com/novalagung/gubrak#Concat) • [Count](https://godoc.org/github.com/novalagung/gubrak#Count) • [Difference](https://godoc.org/github.com/novalagung/gubrak#Difference) • [Drop](https://godoc.org/github.com/novalagung/gubrak#Drop) • [DropRight](https://godoc.org/github.com/novalagung/gubrak#DropRight) • [Each](https://godoc.org/github.com/novalagung/gubrak#Each) • [EachRight](https://godoc.org/github.com/novalagung/gubrak#EachRight) • [Fill](https://godoc.org/github.com/novalagung/gubrak#Fill) • [Filter](https://godoc.org/github.com/novalagung/gubrak#Filter) • [Find](https://godoc.org/github.com/novalagung/gubrak#Find) • [FindIndex](https://godoc.org/github.com/novalagung/gubrak#FindIndex) • [FindLast](https://godoc.org/github.com/novalagung/gubrak#FindLast) • [FindLastIndex](https://godoc.org/github.com/novalagung/gubrak#FindLastIndex) • [First](https://godoc.org/github.com/novalagung/gubrak#First) • [ForEach](https://godoc.org/github.com/novalagung/gubrak#ForEach) • [ForEachRight](https://godoc.org/github.com/novalagung/gubrak#ForEachRight) • [FromPairs](https://godoc.org/github.com/novalagung/gubrak#FromPairs) • [GroupBy](https://godoc.org/github.com/novalagung/gubrak#GroupBy) • [Head](https://godoc.org/github.com/novalagung/gubrak#Head) • [Includes](https://godoc.org/github.com/novalagung/gubrak#Includes) • [IndexOf](https://godoc.org/github.com/novalagung/gubrak#IndexOf) • [Initial](https://godoc.org/github.com/novalagung/gubrak#Initial) • [Intersection](https://godoc.org/github.com/novalagung/gubrak#Intersection) • [IsArray](https://godoc.org/github.com/novalagung/gubrak#IsArray) • [IsBool](https://godoc.org/github.com/novalagung/gubrak#IsBool) • [IsChannel](https://godoc.org/github.com/novalagung/gubrak#IsChannel) • [IsDate](https://godoc.org/github.com/novalagung/gubrak#IsDate) • [IsEmpty](https://godoc.org/github.com/novalagung/gubrak#IsEmpty) • [IsEmptyString](https://godoc.org/github.com/novalagung/gubrak#IsEmptyString) • [IsFloat](https://godoc.org/github.com/novalagung/gubrak#IsFloat) • [IsFunction](https://godoc.org/github.com/novalagung/gubrak#IsFunction) • [IsInt](https://godoc.org/github.com/novalagung/gubrak#IsInt) • [IsMap](https://godoc.org/github.com/novalagung/gubrak#IsMap) • [IsNil](https://godoc.org/github.com/novalagung/gubrak#IsNil) • [IsNumeric](https://godoc.org/github.com/novalagung/gubrak#IsNumeric) • [IsPointer](https://godoc.org/github.com/novalagung/gubrak#IsPointer) • [IsSlice](https://godoc.org/github.com/novalagung/gubrak#IsSlice) • [IsString](https://godoc.org/github.com/novalagung/gubrak#IsString) • [IsStructObject](https://godoc.org/github.com/novalagung/gubrak#IsStructObject) • [IsTrue](https://godoc.org/github.com/novalagung/gubrak#IsTrue) • [IsUint](https://godoc.org/github.com/novalagung/gubrak#IsUint) • [IsZeroNumber](https://godoc.org/github.com/novalagung/gubrak#IsZeroNumber) • [Join](https://godoc.org/github.com/novalagung/gubrak#Join) • [KeyBy](https://godoc.org/github.com/novalagung/gubrak#KeyBy) • [Last](https://godoc.org/github.com/novalagung/gubrak#Last) • [LastIndexOf](https://godoc.org/github.com/novalagung/gubrak#LastIndexOf) • [Map](https://godoc.org/github.com/novalagung/gubrak#Map) • [Now](https://godoc.org/github.com/novalagung/gubrak#Now) • [Nth](https://godoc.org/github.com/novalagung/gubrak#Nth) • [OrderBy](https://godoc.org/github.com/novalagung/gubrak#OrderBy) • [Partition](https://godoc.org/github.com/novalagung/gubrak#Partition) • [Pull](https://godoc.org/github.com/novalagung/gubrak#Pull) • [PullAll](https://godoc.org/github.com/novalagung/gubrak#PullAll) • [PullAt](https://godoc.org/github.com/novalagung/gubrak#PullAt) • [RandomInt](https://godoc.org/github.com/novalagung/gubrak#RandomInt) • [RandomString](https://godoc.org/github.com/novalagung/gubrak#RandomString) • [Reduce](https://godoc.org/github.com/novalagung/gubrak#Reduce) • [Reject](https://godoc.org/github.com/novalagung/gubrak#Reject) • [Remove](https://godoc.org/github.com/novalagung/gubrak#Remove) • [Reverse](https://godoc.org/github.com/novalagung/gubrak#Reverse) • [Sample](https://godoc.org/github.com/novalagung/gubrak#Sample) • [SampleSize](https://godoc.org/github.com/novalagung/gubrak#SampleSize) • [Shuffle](https://godoc.org/github.com/novalagung/gubrak#Shuffle) • [Size](https://godoc.org/github.com/novalagung/gubrak#Size) • [SortBy](https://godoc.org/github.com/novalagung/gubrak#SortBy) • [Tail](https://godoc.org/github.com/novalagung/gubrak#Tail) • [Take](https://godoc.org/github.com/novalagung/gubrak#Take) • [TakeRight](https://godoc.org/github.com/novalagung/gubrak#TakeRight) • [Union](https://godoc.org/github.com/novalagung/gubrak#Union) • [Uniq](https://godoc.org/github.com/novalagung/gubrak#Uniq) • [Without](https://godoc.org/github.com/novalagung/gubrak#Without)

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
