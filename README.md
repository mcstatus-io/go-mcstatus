# go-mcstatus
The official Go library for interacting with the [mcstatus.io](https://mcstatus.io) API.

## Getting Started

Firstly, you will need to install the library. Open your terminal/command line in your workspace and run the following command.

```
go get github.com/mcstatus-io/go-mcstatus
```

## Usage

### Java Status

```go
package main

import (
    "fmt"

    "github.com/mcstatus-io/go-mcstatus"
)

func main() {
    resp, err := mcstatus.GetJavaStatus("demo.mcstatus.io", 25565)

    if err != nil {
        panic(err)
    }

    fmt.Printf("%+v\n", resp)
}
```

### Bedrock Status

```go
package main

import (
    "fmt"
    
    "github.com/mcstatus-io/go-mcstatus"
)

func main() {
    resp, err := mcstatus.GetBedrockStatus("demo.mcstatus.io", 19132)

    if err != nil {
        panic(err)
    }

    fmt.Printf("%+v\n", resp)
}
```

### Java Widget

```go
package main

import "github.com/mcstatus-io/go-mcstatus"

func main() {
    img, err := mcstatus.GetJavaWidget("demo.mcstatus.io", 25565)

    if err != nil {
        panic(err)
    }

    // ...
}
```

### Icon

```go
package main

import "github.com/mcstatus-io/go-mcstatus"

func main() {
    img, err := mcstatus.GetIcon("demo.mcstatus.io", 25565)

    if err != nil {
        panic(err)
    }

    // ...
}
```

### Default Icon

```go
package main

import "github.com/mcstatus-io/go-mcstatus"

func main() {
    img := mcstatus.GetDefaultIcon()

    // ...
}
```

## License
[MIT License](https://github.com/mcstatus-io/go-mcstatus/blob/main/LICENSE)