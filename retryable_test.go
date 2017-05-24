package errorutil

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	pkgerrors "github.com/pkg/errors"
)

type retryable bool

func (err retryable) Error() string {
	return fmt.Sprintf("%v", err)
}

func (err retryable) Retryable() bool {
	return bool(err)
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		err  error
		want bool
	}{
		{nil, false},
		{RetryableError(nil), false},
		{pkgerrors.Wrap(RetryableError(nil), "bar"), false},

		{errors.New("foo"), false},
		{pkgerrors.New("foo"), false},
		{pkgerrors.Wrap(errors.New("foo"), "bar"), false},

		{RetryableError(errors.New("foo")), true},
		{RetryableError(pkgerrors.New("foo")), true},
		{pkgerrors.Wrap(RetryableError(errors.New("foo")), "bar"), true},

		{httpError(http.StatusBadGateway), true},
		{pkgerrors.Wrap(httpError(http.StatusBadGateway), "bar"), true},
		{httpError(http.StatusInternalServerError), true},
		{pkgerrors.Wrap(httpError(http.StatusInternalServerError), "bar"), true},
		{httpError(http.StatusGatewayTimeout), true},
		{pkgerrors.Wrap(httpError(http.StatusGatewayTimeout), "bar"), true},
		{httpError(429), true},
		{pkgerrors.Wrap(httpError(429), "bar"), true},
		{httpError(http.StatusNotFound), false},
		{pkgerrors.Wrap(httpError(http.StatusNotFound), "bar"), false},
		{NotFoundError(errors.New("foo")), false},
		{httpError(http.StatusBadRequest), false},
		{pkgerrors.Wrap(httpError(http.StatusBadRequest), "bar"), false},
		{InvalidError(errors.New("foo")), false},
		{httpError(http.StatusForbidden), false},
		{pkgerrors.Wrap(httpError(http.StatusForbidden), "bar"), false},
		{ForbiddenError(errors.New("foo")), false},

		{retryable(false), false},
		{retryable(true), true},
	}
	for _, tt := range tests {
		got := IsRetryable(tt.err)
		if got != tt.want {
			t.Errorf("IsRetryable(%q): got: %v, want %v", tt.err, got, tt.want)
		}
	}
}

func ExampleIsRetryable() {
	err := errors.New("some error")
	IsRetryable(err) // will return false
	err = RetryableError(err)
	IsRetryable(err) // will return true
}

func TestIsNotRetryable(t *testing.T) {
	tests := []struct {
		err  error
		want bool
	}{
		{nil, false},
		{RetryableError(nil), false},
		{pkgerrors.Wrap(RetryableError(nil), "bar"), false},

		{errors.New("foo"), false},
		{pkgerrors.New("foo"), false},
		{pkgerrors.Wrap(errors.New("foo"), "bar"), false},

		{RetryableError(errors.New("foo")), false},
		{RetryableError(pkgerrors.New("foo")), false},
		{pkgerrors.Wrap(RetryableError(errors.New("foo")), "bar"), false},

		{retryable(false), true},
		{retryable(true), false},
	}
	for _, tt := range tests {
		got := IsNotRetryable(tt.err)
		if got != tt.want {
			t.Errorf("IsRetryable(%q): got: %v, want %v", tt.err, got, tt.want)
		}
	}
}

func ExampleIsNotRetryable() {
	err := errors.New("some error")
	IsNotRetryable(err) // will return false
	err = RetryableError(err)
	IsNotRetryable(err) // will return false
}
