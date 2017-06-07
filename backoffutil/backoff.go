// Package backoffutil provides a wrapper above github.com/cenk/backoff.Retry
// that checks the error returned and only retries retryable errors.
//
// For the sake of simplicity, the backoff strategy is the default exponential backoff.
package backoffutil

import (
	"github.com/cenk/backoff"
	"github.com/objenious/errorutil"
)

// Retry does exponential backoff.
// Backoff will trigger if an error is returned, implements Retryabler AND the error is retryable.
func Retry(fn func() error) error {
	return RetryWithCustomBackOff(fn, backoff.NewExponentialBackOff())
}

// RetryWithCustomExponentialBackOff does exponential backoff with a given BackOff instance.
// Backoff will trigger if an error is returned, implements Retryabler AND the error is retryable.
func RetryWithCustomBackOff(fn func() error, bo backoff.BackOff) error {
	var finalerr error
	err := backoff.Retry(func() error {
		finalerr = fn()
		if errorutil.IsRetryable(finalerr) {
			return finalerr
		}
		return nil
	}, bo)

	if err != nil {
		return err
	}
	return finalerr
}
