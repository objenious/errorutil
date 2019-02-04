package errorutil

import (
	"github.com/pkg/errors"
	"time"
)

// Delayer defines an error that
type Delayer interface {
	Delay() time.Duration
}

// Delay return the delay duration of a DelayedError (i.e. implements Delayer).
// If the error is nil or does not implement Delayer or the delay is not a positive value, 0 is returned.
func Delay(err error) time.Duration {
	type causer interface {
		Cause() error
	}

	for err != nil {
		if delay, ok := err.(Delayer); ok {
			return delay.Delay()
		}
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return 0
}

// DelayedError set error delay duration.  returns nil if the error is nil.
func WithDelay(err error, duration time.Duration) error {
	if err == nil {
		return nil
	}
	return &delayedError{err, duration}
}

// NewDelayedError returns a delayed error that formats as the given text and duration.
func NewDelayedError(text string, duration time.Duration) error {
	return WithDelay(errors.New(text), duration)
}

type delayedError struct {
	error
	duration time.Duration
}

func (err *delayedError) Error() string {
	return err.error.Error()
}

func (err *delayedError) Delay() time.Duration {
	return err.duration
}

func (err *delayedError) Cause() error {
	return err.error
}
