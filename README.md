# TypedEnv

Strongly typed environment variable management for Go.

## Features

- Errors for missing environment variables
- Errors for type mismatches and parsing failures
- Variable values are not exposed in logs or output
- No unsafe
- No panics

## Installation

``` bash
go get github.com/saas-craft/typedenv
```

## Usage

``` go
package main

import (
    "fmt"
    "log"

    "github.com/saas-craft/typedenv"
)

func main() {
    type config struct {
        AppHost string `env:"HOST"`
        AppPort int    `env:"PORT"`
    }

    if cfg, err := typedenv.Load[config](); err != nil {
        log.Fatalf("load config: %w", err)
    }

    fmt.Println(cfg)
}
```

## License

SaaS Craft TypedEnv is licensed under the MIT License - see [LICENSE](LICENSE) for details.
