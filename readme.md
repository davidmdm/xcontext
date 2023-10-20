# xcontext

The `xcontext` package is a Go library that enhances the capabilities of the standard `context` package by adding signal-based context cancellation. This functionality allows you to cancel a context when specific signals are received, such as interrupt signals (e.g., SIGINT) or termination signals (e.g., SIGTERM). The package is particularly useful in applications where graceful shutdown or cleanup is required in response to system signals.

## Installation

You can easily add the `xcontext` package to your Go project using the following go get command:

```bash
go get github.com/davidmdm/xcontext
```

## Usage

### `WithSignalCancelation`

Here is a simple example of how to use the `xcontext` package to gracefully handle SIGINT (Ctrl+C) and SIGTERM (termination signal) in your application:

```go
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "github.com/davidmdm/xcontext"
)

func main() {
    // Create a context that cancels on SIGINT and SIGTERM signals
    ctx, cancel := xcontext.WithSignalCancelation(context.Background(), sycall.SIGINT, syscall.SIGTERM)
    defer canel()

    // Will be done once context is canceled by a SIGINT
    <-ctx.Done()

    context.Err()            // context.Canceled
    context.Cause(ctx)       // xcontext.SignalContextError -> context canceled: signal received: interrupt
    xcontext.SignalCause(ctx) // syscall.SIGINT (the received signal value)
}
```

In this example, the application uses the `WithSignalCancelation` function to create a context that cancels when either SIGINT or SIGTERM is received.

### `SignalCancelError`

The `SignalCancelError` is a custom error type defined in this package. When a context is canceled due to a signal, this error is used as the context cause, providing information about the signal received. You can use this error type to handle signal-related errors in your code.

### `SignalCause`

The `SignalCause` function allows you to retrieve the signal that caused a context to be canceled. If the context was canceled due to a signal, it returns the signal; otherwise, it returns `nil`. This can be helpful if you want to log or handle the specific signal that triggered the context cancellation.

## Contribution

If you want to contribute to this package or report any issues, please visit the GitHub repository at [https://github.com/davidmdm/xcontext](https://github.com/davidmdm/xcontext).

## License

This package is open source and is provided under the [MIT License](https://github.com/davidmdm/xcontext/blob/main/LICENSE).
