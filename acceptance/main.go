package main

import (
	"context"
	"fmt"
	"syscall"
	"time"

	"github.com/davidmdm/env"
	"github.com/davidmdm/xcontext"
)

type Config struct{}

func main() {
	var (
		rootTimeout time.Duration
		useCancel   bool
	)

	env.Var(&rootTimeout, "ROOT_TIMEOUT")
	env.Var(&useCancel, "USE_CANCEL")

	env.MustParse()

	ctx, cancel := func() (context.Context, context.CancelFunc) {
		if rootTimeout == 0 {
			return context.Background(), func() {}
		}
		return context.WithTimeout(context.Background(), rootTimeout)
	}()

	defer cancel()

	ctx, cancel = xcontext.WithSignalCancelation(ctx, syscall.SIGINT)
	defer cancel()

	if useCancel {
		cancel()
	}

	<-ctx.Done()

	fmt.Println(context.Cause(ctx))
}
