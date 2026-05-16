# TypedEnv [![CI](https://github.com/saas-craft/typedenv/actions/workflows/ci.yml/badge.svg)](https://github.com/saas-craft/typedenv/actions/workflows/ci.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/saas-craft/typedenv)](https://goreportcard.com/report/github.com/saas-craft/typedenv)

Strongly typed environment variable management for Go.

Configure your requirements using a Go struct. TypedEnv validates required environment variables and returns explicit errors for missing or invalid values. .env files are intentionally unsupported to reduce accidental secret exposure, and to keep configuration sourcing explicit.

## Features

- Errors for missing environment variables
- Errors for type mismatches and parsing failures
- No variable values in errors or logs
- No iteration of OS variables; look only at what's necessary
- No unsafe usage
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

    cfg, err := typedenv.Load[config]()
    if err != nil {
        log.Fatalf("load config: %v", err)
    }

    fmt.Println(cfg)
}
```

## License

SaaS Craft TypedEnv is licensed under the MIT License - see [LICENSE](LICENSE) for details.
