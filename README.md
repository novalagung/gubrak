# Gubrak v2

Gubrak is Go functional utility library with syntactic sugar. It's like lodash, but for Go Programming language.

[![Go Report Card](https://goreportcard.com/badge/github.com/novalagung/gubrak?nocache=1)](https://goreportcard.com/report/github.com/novalagung/gubrak?nocache=1)
[![Build Status](https://travis-ci.org/novalagung/gubrak.svg?branch=master)](https://travis-ci.org/novalagung/gubrak)
[![Coverage Status](https://coveralls.io/repos/github/novalagung/gubrak/badge.svg?branch=master)](https://coveralls.io/github/novalagung/gubrak?branch=master)

## Installation

The latest version of gubrak is v2. Here are the available method to get this library.

- Using `go get` from github, for **Go Mod**-based project:

    ```bash
    go get -u github.com/novalagung/gubrak/v2
    ```

- Using `go get` from github, for `$GOPATH`-based project:

    ```bash
    go get -u github.com/novalagung/gubrak
    ```

## Usage

Since this library is developed as go module, the versioning system used is the one from Go.

- For **Go Mod**-based project:

    ```go
    import gubrak "github.com/novalagung/gubrak/v2"
    ```

- For `$GOPATH`-based project:

    ```go
    import "github.com/novalagung/gubrak"
    ```

## Documentation

 - [Pkg Dev](https://pkg.go.dev/github.com/novalagung/gubrak/v2)
 - [Godoc](https://godoc.org/github.com/novalagung/gubrak)

## Hello World Example

![Lodash for Golang](https://i.imgur.com/bvT1gVM.jpg)

## APIs

Below are the list of available functions on gubrak:

[Chunk](https://pkg.go.dev/github.com/novalagung/gubrak#Chunk) • [Compact](https://pkg.go.dev/github.com/novalagung/gubrak#Compact) • [Concat](https://pkg.go.dev/github.com/novalagung/gubrak#Concat) • [Count](https://pkg.go.dev/github.com/novalagung/gubrak#Count) • [Difference](https://pkg.go.dev/github.com/novalagung/gubrak#Difference) • [Drop](https://pkg.go.dev/github.com/novalagung/gubrak#Drop) • [DropRight](https://pkg.go.dev/github.com/novalagung/gubrak#DropRight) • [Each](https://pkg.go.dev/github.com/novalagung/gubrak#Each) • [EachRight](https://pkg.go.dev/github.com/novalagung/gubrak#EachRight) • [Fill](https://pkg.go.dev/github.com/novalagung/gubrak#Fill) • [Filter](https://pkg.go.dev/github.com/novalagung/gubrak#Filter) • [Find](https://pkg.go.dev/github.com/novalagung/gubrak#Find) • [FindIndex](https://pkg.go.dev/github.com/novalagung/gubrak#FindIndex) • [FindLast](https://pkg.go.dev/github.com/novalagung/gubrak#FindLast) • [FindLastIndex](https://pkg.go.dev/github.com/novalagung/gubrak#FindLastIndex) • [First](https://pkg.go.dev/github.com/novalagung/gubrak#First) • [ForEach](https://pkg.go.dev/github.com/novalagung/gubrak#ForEach) • [ForEachRight](https://pkg.go.dev/github.com/novalagung/gubrak#ForEachRight) • [FromPairs](https://pkg.go.dev/github.com/novalagung/gubrak#FromPairs) • [GroupBy](https://pkg.go.dev/github.com/novalagung/gubrak#GroupBy) • [Head](https://pkg.go.dev/github.com/novalagung/gubrak#Head) • [Includes](https://pkg.go.dev/github.com/novalagung/gubrak#Includes) • [IndexOf](https://pkg.go.dev/github.com/novalagung/gubrak#IndexOf) • [Initial](https://pkg.go.dev/github.com/novalagung/gubrak#Initial) • [Intersection](https://pkg.go.dev/github.com/novalagung/gubrak#Intersection) • [IsArray](https://pkg.go.dev/github.com/novalagung/gubrak#IsArray) • [IsBool](https://pkg.go.dev/github.com/novalagung/gubrak#IsBool) • [IsChannel](https://pkg.go.dev/github.com/novalagung/gubrak#IsChannel) • [IsDate](https://pkg.go.dev/github.com/novalagung/gubrak#IsDate) • [IsEmpty](https://pkg.go.dev/github.com/novalagung/gubrak#IsEmpty) • [IsEmptyString](https://pkg.go.dev/github.com/novalagung/gubrak#IsEmptyString) • [IsFloat](https://pkg.go.dev/github.com/novalagung/gubrak#IsFloat) • [IsFunction](https://pkg.go.dev/github.com/novalagung/gubrak#IsFunction) • [IsInt](https://pkg.go.dev/github.com/novalagung/gubrak#IsInt) • [IsMap](https://pkg.go.dev/github.com/novalagung/gubrak#IsMap) • [IsNil](https://pkg.go.dev/github.com/novalagung/gubrak#IsNil) • [IsNumeric](https://pkg.go.dev/github.com/novalagung/gubrak#IsNumeric) • [IsPointer](https://pkg.go.dev/github.com/novalagung/gubrak#IsPointer) • [IsSlice](https://pkg.go.dev/github.com/novalagung/gubrak#IsSlice) • [IsString](https://pkg.go.dev/github.com/novalagung/gubrak#IsString) • [IsStructObject](https://pkg.go.dev/github.com/novalagung/gubrak#IsStructObject) • [IsTrue](https://pkg.go.dev/github.com/novalagung/gubrak#IsTrue) • [IsUint](https://pkg.go.dev/github.com/novalagung/gubrak#IsUint) • [IsZeroNumber](https://pkg.go.dev/github.com/novalagung/gubrak#IsZeroNumber) • [Join](https://pkg.go.dev/github.com/novalagung/gubrak#Join) • [KeyBy](https://pkg.go.dev/github.com/novalagung/gubrak#KeyBy) • [Last](https://pkg.go.dev/github.com/novalagung/gubrak#Last) • [LastIndexOf](https://pkg.go.dev/github.com/novalagung/gubrak#LastIndexOf) • [Map](https://pkg.go.dev/github.com/novalagung/gubrak#Map) • [Now](https://pkg.go.dev/github.com/novalagung/gubrak#Now) • [Nth](https://pkg.go.dev/github.com/novalagung/gubrak#Nth) • [OrderBy](https://pkg.go.dev/github.com/novalagung/gubrak#OrderBy) • [Partition](https://pkg.go.dev/github.com/novalagung/gubrak#Partition) • [Pull](https://pkg.go.dev/github.com/novalagung/gubrak#Pull) • [PullAll](https://pkg.go.dev/github.com/novalagung/gubrak#PullAll) • [PullAt](https://pkg.go.dev/github.com/novalagung/gubrak#PullAt) • [RandomInt](https://pkg.go.dev/github.com/novalagung/gubrak#RandomInt) • [RandomString](https://pkg.go.dev/github.com/novalagung/gubrak#RandomString) • [Reduce](https://pkg.go.dev/github.com/novalagung/gubrak#Reduce) • [Reject](https://pkg.go.dev/github.com/novalagung/gubrak#Reject) • [Remove](https://pkg.go.dev/github.com/novalagung/gubrak#Remove) • [Reverse](https://pkg.go.dev/github.com/novalagung/gubrak#Reverse) • [Sample](https://pkg.go.dev/github.com/novalagung/gubrak#Sample) • [SampleSize](https://pkg.go.dev/github.com/novalagung/gubrak#SampleSize) • [Shuffle](https://pkg.go.dev/github.com/novalagung/gubrak#Shuffle) • [Size](https://pkg.go.dev/github.com/novalagung/gubrak#Size) • [SortBy](https://pkg.go.dev/github.com/novalagung/gubrak#SortBy) • [Tail](https://pkg.go.dev/github.com/novalagung/gubrak#Tail) • [Take](https://pkg.go.dev/github.com/novalagung/gubrak#Take) • [TakeRight](https://pkg.go.dev/github.com/novalagung/gubrak#TakeRight) • [Union](https://pkg.go.dev/github.com/novalagung/gubrak#Union) • [Uniq](https://pkg.go.dev/github.com/novalagung/gubrak#Uniq) • [Without](https://pkg.go.dev/github.com/novalagung/gubrak#Without)

## Test

```bash
go test -cover -race -v ./... 
```

## Contribution

Fork ➜ Create branch ➜ Commit ➜ Push ➜ Pull Requests

## License

MIT License
