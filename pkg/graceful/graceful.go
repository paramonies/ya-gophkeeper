// Package graceful contains API to implement application graceful shutdown.
// Application starts listening for SIGINT or SIGTERM signals and handles them properly.
package graceful

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	ShutdownTimeout = 10 * time.Second
)

// ShutdownFunc is a callback-type for registering callbacks before application shutdown.
type ShutdownFunc func() error

var (
	// ErrTimeoutExceeded is returned when the application fails to shut down for a given period of time.
	ErrTimeoutExceeded = errors.New("failed to perform graceful shutdown: timed out")

	// ErrForceShutdown is returned when the user or operating system is sending SIGINT or SIGTERM
	// for the application being is graceful-shutdown state.
	ErrForceShutdown = errors.New("failed to perform graceful shutdown: force shut downed")
)

var (
	handler   *shutdownHandler
	execOnErr func(error)
)

func init() {
	setupHandler()
}

func setupHandler() {
	notify := make(chan os.Signal, 1)
	signal.Notify(notify, syscall.SIGINT, syscall.SIGTERM)

	handler = newHandler(notify)

	execOnErr = func(err error) {
		log.Printf("shutdown callback error: %v", err)
	}
}

// AddCallback registers a callback for execution before shutdown.
func AddCallback(fn ShutdownFunc) {
	handler.add(fn)
}

// ExecOnError executes the given handler
// when shutdown callback returns any error.
func ExecOnError(cb func(err error)) {
	execOnErr = cb
}

// WaitShutdown waits for application shutdown.
//
// If the user or operating system interrupts the graceful shutdown,
// ErrForceShutdown is returned.
//
// If application fails to shut down for a given period of time,
// ErrTimeoutExceeded is returned.
func WaitShutdown() error {
	<-handler.C

	notify := make(chan os.Signal, 1)
	signal.Notify(notify, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)

		callbacks := handler.get()
		for i := len(callbacks) - 1; i >= 0; i-- {
			err := callbacks[i]()
			if err != nil && execOnErr != nil {
				execOnErr(err)
			}
		}
	}()

	select {
	case <-done:
		return nil
	case <-notify:
		return ErrForceShutdown
	case <-ctx.Done():
		return ErrTimeoutExceeded
	}
}

// ShutdownNow sends interrupt signal to initiate graceful shutdown.
func ShutdownNow() {
	handler.C <- os.Interrupt
}
