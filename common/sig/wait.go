package sig

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

var (
	// syscall.SIGINT, syscall.SIGTERM
	DefaultSigs = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	DefaultFunc = func() error { return nil }
)

func Wait(sigs []os.Signal, logWriter io.Writer, callbacks ...func() error) {
	<-Chan(sigs, logWriter, callbacks...)
}

func Chan(sigs []os.Signal, logWriter io.Writer, callbacks ...func() error) <-chan bool {
	if sigs == nil {
		sigs = DefaultSigs
	}
	if logWriter == nil {
		logWriter = os.Stdout
	}
	var (
		sig  = make(chan os.Signal, 1)
		done = make(chan bool, 1)
	)
	signal.Notify(sig, sigs...)
	go func() {
		s := <-sig

		errs := make([]error, 0, len(callbacks))
		for _, fn := range callbacks {
			errs = append(errs, fn())
		}
		fmt.Fprintln(logWriter, "input signal: ", s.String(), errors.Join(errs...))

		done <- true
	}()
	return done
}
