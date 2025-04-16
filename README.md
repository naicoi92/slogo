# SLogo

SLogo is a lightweight, customizable formatter and handler for the Go standard library's `log/slog` package. It enhances logging capabilities by providing structured logging with custom formatters and multiple handler support.

[![Go Reference](https://pkg.go.dev/badge/github.com/naicoi92/slogo.svg)](https://pkg.go.dev/github.com/naicoi92/slogo)
[![Go Report Card](https://goreportcard.com/badge/github.com/naicoi92/slogo)](https://goreportcard.com/report/github.com/naicoi92/slogo)

## Features

- **Custom Formatters**: Format struct fields, errors, and other complex data types in a human-readable format
- **Multiple Handler Support**: Route logs to multiple destinations using the slogmulti fanout capability
- **Configurable Options**: Easily customize logging behavior with functional options
- **Enhanced Error Reporting**: Capture error details including type and message
- **Struct Field Handling**: Control how struct fields are logged, including field redaction
- **JSON Tag Support**: Respects JSON field tags for naming consistency

## Installation

```bash
go get github.com/naicoi92/slogo
```

## Quick Start

```go
package main

import (
    "errors"
    "log/slog"
    "os"
    
    "github.com/naicoi92/slogo"
)

func main() {
    // Create a new logger with default settings
    logger := slog.New(
        slogo.NewHandler(
            os.Stdout,
            &slog.HandlerOptions{
                Level:     slog.LevelInfo,
                AddSource: true,
            },
        ),
    )
    
    // Set as default logger
    slog.SetDefault(logger)
    
    // Log messages
    slog.Info("Hello from SLogo")
    
    // Log with attributes
    slog.Info("User login", "user_id", 12345, "action", "login")
    
    // Log structs
    type User struct {
        ID        int    `json:"id"`
        Username  string `json:"username"`
        Password  string `json:"password" slog:"restrict"` // This field will be redacted
        Email     string `json:"email"`
    }
    
    user := User{
        ID:       1,
        Username: "johndoe",
        Password: "secret123",
        Email:    "john@example.com",
    }
    
    slog.Info("User details", "user", user)
    
    // Log errors with type and message information
    if err := errors.New("an example error"); err != nil {
        slog.Error("Operation failed", "error", err)
    }
}
```

## Complex Struct Example

SLogo excels at formatting complex nested structs:

```go
type User struct {
    ID        int        `json:"id"`
    Username  string     `json:"username"`
    Password  string     `json:"password" slog:"restrict"` // This field will be redacted
    Email     string     `json:"email"`
    IpAddress netip.Addr `json:"ip_address"`
    Members   []*User    `json:"members"`
}

// Create a user with nested members
ip, _ := netip.ParseAddr("10.0.0.1")
user := User{
    ID:        1,
    Username:  "johndoe",
    Password:  "secret123",
    Email:     "john@example.com",
    IpAddress: ip,
    Members: []*User{
        {
            ID:       2,
            Username: "janedoe",
        },
        {
            ID:        3,
            Username:  "alice",
            IpAddress: ip,
        },
    },
}

slog.Info("User details", "user", user)
```

## Configuration Options

SLogo provides several configuration options through functional options:

### WithSlogor

Use the underlying slogor handler:

```go
logger := slog.New(slogo.NewHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
}, slogo.WithSlogor()))
```

### WithSlogHandler

Add a custom slog.Handler:

```go
customHandler := // your custom handler
logger := slog.New(slogo.NewHandler(os.Stdout, &slog.HandlerOptions{}, 
    slogo.WithSlogHandler(customHandler)))
```

### WithFormatter

Add a custom formatter:

```go
customFormatter := // your custom formatter
logger := slog.New(slogo.NewHandler(os.Stdout, &slog.HandlerOptions{}, 
    slogo.WithFormatter(customFormatter)))
```

## Built-in Formatters

### FormatStruct

The `FormatStruct` formatter handles struct data types and provides special handling for struct fields:

- Uses JSON field tags for field names if available
- Recursively formats nested structs and slices
- Skips empty values
- Respects `slog` struct tags:
  - `-`: Skip this field entirely
  - `restrict`: Redact the field value (shows "[REDACTED]")

Example of formatting a struct with redacted fields:

```
[id=1, username=johndoe, password=[REDACTED], email=john@example.com, ip_address=10.0.0.1, members=[[id=2, username=janedoe], [id=3, username=alice, ip_address=10.0.0.1]]]
```

### FormatError

The `FormatError` formatter enhances error logging by capturing:

- Error message
- Error type

Example of formatted error output:

```
[Message=an example error, Type=*errors.errorString]
```

## Dependencies

SLogo depends on the following packages:

- [github.com/samber/slog-formatter](https://github.com/samber/slog-formatter)
- [github.com/samber/slog-multi](https://github.com/samber/slog-multi)
- [gitlab.com/greyxor/slogor](https://gitlab.com/greyxor/slogor)

## Complete Example

Check out the [example](./example/main.go) directory for a complete working example of SLogo.

## License

This project is licensed under the terms found in the [LICENSE](./LICENSE) file.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request