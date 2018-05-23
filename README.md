# GUBRAK

Golang utility library with syntactic sugar.

[![Go Report Card](https://goreportcard.com/badge/github.com/novalagung/gubrak?nocache=1)](https://goreportcard.com/report/github.com/novalagung/gubrak?nocache=1)
[![Build Status](https://travis-ci.org/novalagung/gubrak.svg?branch=master)](https://travis-ci.org/novalagung/gubrak)
[![Coverage Status](https://coveralls.io/repos/github/novalagung/gubrak/badge.svg?branch=master)](https://coveralls.io/github/novalagung/gubrak?branch=master)
<!-- [![codecov](https://codecov.io/gh/novalagung/gubrak/branch/master/graph/badge.svg)](https://codecov.io/gh/novalagung/gubrak) -->
<!-- [![cover.run](https://cover.run/go/https:/github.com/novalagung/gubrak.svg?style=flat&tag=golang-1.10)](https://cover.run/go?tag=golang-1.10&repo=https%3A%2Fgithub.com%2Fnovalagung%2Fgubrak) -->


Gubrak is yet another utility library for Golang, inspired from lodash. Currently we have around 46 reusable functions available, and it'll keep increasing.

### Installation

```go
go get -u github.com/novalagung/gubrak
```

### Documentation

See [official documentation page](https://gubrak.github.io/#/documentation).

[Godoc documentation page](http://godoc.org/github.com/novalagung/gubrak) is still in progress.

### Hello World Example

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

### Contribution

Fork ➜ Create branch ➜ Commit ➜ Push ➜ Pull Requests

### License

MIT License
