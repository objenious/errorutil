package errorutil

import (
	"errors"
	"fmt"
	"net/http"
)

// Retryabler defines an error that may be temporary. A function returning a retryable error should be executed again.
type Retryabler interface {
	Retryable() bool
}

// IsRetryable checks if an error is retryable (i.e. implements Retryabler and Retryable returns true).
//
// If the error is nil or does not implement Retryabler, false is returned.
func IsRetryable(err error) bool {
	type causer interface {
		Cause() error
	}

	for err != nil {
		if retry, ok := err.(Retryabler); ok {
			return retry.Retryable()
		}
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return false
}

// IsNotRetryable checks if an error is explicitly marked as not retryable (i.e. implements Retryabler and Retryable returns false).
//
// If the error is nil or does not implement Retryabler, false is returned.
func IsNotRetryable(err error) bool {
	type causer interface {
		Cause() error
	}

	for err != nil {
		if retry, ok := err.(Retryabler); ok {
			return !retry.Retryable()
		}
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return false
}

// RetryableError marks an error as retryable. It returns nil if the error is nil.
func RetryableError(err error) error {
	if err == nil {
		return nil
	}
	return &retryableError{err: err}
}

// NotRetryableError marks an error as NOT retryable. It returns nil if the error is nil.
func NotRetryableError(err error) error {
	if err == nil {
		return nil
	}
	return &notRetryableError{err: err}
}

// NewRetryableError returns a retryable error that formats as the given text.
func NewRetryableError(text string) error {
	return RetryableError(errors.New(text))
}

// NewRetryableErrorf formats according to a format specifier and returns the string
// as a value that satisfies a retryable error.
func NewRetryableErrorf(format string, args ...interface{}) error {
	return RetryableError(fmt.Errorf(format, args...))
}

type retryableError struct {
	err error
}

func (err *retryableError) Error() string {
	return err.err.Error()
}

func (err *retryableError) Retryable() bool {
	return true
}

func (err *retryableError) HTTPStatusCode() int {
	return http.StatusInternalServerError
}

func (err *retryableError) Cause() error {
	return err.err
}

type notRetryableError struct {
	err error
}

func (err *notRetryableError) Error() string {
	return err.err.Error()
}

func (err *notRetryableError) Retryable() bool {
	return false
}

func (err *notRetryableError) HTTPStatusCode() int {
	return http.StatusInternalServerError
}

func (err *notRetryableError) Cause() error {
	return err.err
}
