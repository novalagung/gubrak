# GUBRAK

Golang utility library with syntactic sugar.

Gubrak is yet another utility library for Golang, inspired from lodash. Currently we have around 46 reusable functions available, and it'll keep increasing.

### Installation

```go
go get -u github.com/novalagung/gubrak
```

### Documentation

See [documentation page](https://gubrak.github.io/#/documentation).

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

Test files are available inside `/test` folder.

### License

MIT License
