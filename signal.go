package xcontext

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
)

func WithSignalCancelation(parent context.Context, signals ...os.Signal) (ctx context.Context, cancel context.CancelFunc) {
	var (
		signalCh = make(chan os.Signal, 1)
		done     = make(chan struct{})
		stop     = make(chan struct{})
	)

	signal.Notify(signalCh, signals...)

	ctx, cancelCause := context.WithCancelCause(parent)

	go func() {
		defer close(done)
		defer signal.Stop(signalCh)

		select {
		case sig := <-signalCh:
			cancelCause(SignalCancelError{sig})
			return
		case <-parent.Done():
		case <-stop:
		}

		cancelCause(nil)
	}()

	var once sync.Once

	cancel = func() {
		once.Do(func() { close(stop) })
		<-done
	}

	return ctx, cancel
}

type SignalCancelError struct {
	Signal os.Signal
}

func (err SignalCancelError) Error() string {
	return fmt.Sprintf("%v: received signal: %s", context.Canceled, err.Signal)
}

func (SignalCancelError) Unwrap() error {
	return context.Canceled
}

func SignalCause(ctx context.Context) os.Signal {
	if sigErr := (SignalCancelError{}); errors.As(context.Cause(ctx), &sigErr) {
		return sigErr.Signal
	}
	return nil
}
