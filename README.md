# TypedEnv [![CI](https://github.com/saas-craft/typedenv/actions/workflows/ci.yml/badge.svg)](https://github.com/saas-craft/typedenv/actions/workflows/ci.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/saas-craft/typedenv)](https://goreportcard.com/report/github.com/saas-craft/typedenv) [![Go Reference](https://pkg.go.dev/badge/github.com/saas-craft/typedenv.svg)](https://pkg.go.dev/github.com/saas-craft/typedenv)

Type-safe environment configuration for Go.

```go
type Config struct {
    Host     string        `env:"HOST"`
    Port     int           `env:"PORT"`
    Timeout  time.Duration `env:"TIMEOUT"`
    LogLevel slog.Level    `env:"LOG_LEVEL,default=debug"`
}

cfg, err := typedenv.Load[Config]()
```

Unparseable and missing key values are returned in one error.

## Features

- Optional default value for keys
- Looks up only the keys you declare
- Keeps raw values out of errors
- No use of the unsafe package
- No panics
- No .env support promotes explicit sourcing

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
    "os"
    "log/slog"
    "time"

    "github.com/saas-craft/typedenv"
)

func main() {
    type config struct {
        Host     string        `env:"HOST"`
        Port     int           `env:"PORT"`
        Timeout  time.Duration `env:"TIMEOUT"`
        LogLevel slog.Level    `env:"LOG_LEVEL,default=debug"`
    }

    os.Setenv("HOST", "localhost")
    os.Setenv("PORT", "8080")
    os.Setenv("TIMEOUT", "1s")

    cfg, err := typedenv.Load[config]()
    if err != nil {
        log.Fatalf("load config: %v", err)
    }

    fmt.Printf("%#v\n", cfg)
    // Output: typedenv.config{Host:"localhost", Port:8080, Timeout:1000000000, LogLevel:-4}
}
```

## Supported Types

| Go Type | Example value |
| --- | --- |
| `string` | `hello` |
| `bool` | `true`, `false`, `1`, `0` |
| `int`, `int8`, `int16`, `int32`, `int64` | `-42` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `42` |
| `float32`, `float64` | `3.14` |
| `time.Duration` | `1h30m`, `500ms`, `2s` |
| `url.URL` | `https://saascraft.com/v1` |
| `encoding.TextUnmarshaler` (e.g. `slog.Level`) | `debug` |

Untagged fields are left at their zero value.

## Works Well With

- SaasCraft [Secret](https://pkg.go.dev/github.com/saas-craft/secret), a generic wrapper that hides values from default formatting, logging, and serialization

## Constraints

- No support for named time.Duration wrapper types, which can't be distinguished from integers

## License

SaasCraft TypedEnv is licensed under the MIT License - see [LICENSE](LICENSE) for details.
