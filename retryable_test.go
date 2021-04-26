package errorutil

import (
	"errors"
	"net/http"
	"testing"

	oerrors "github.com/objenious/errors"
)

type retryable bool

func (err retryable) Error() string {
	if err {
		return "true"
	}
	return "false"
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
		{NotRetryableError(nil), false},
		{oerrors.Wrap(RetryableError(nil), "bar"), false},

		{errors.New("foo"), false},
		{oerrors.New("foo"), false},
		{oerrors.Wrap(errors.New("foo"), "bar"), false},

		{RetryableError(errors.New("foo")), true},
		{RetryableError(oerrors.New("foo")), true},
		{oerrors.Wrap(RetryableError(errors.New("foo")), "bar"), true},

		{RetryableError(NotRetryableError(errors.New("foo"))), true},
		{NotRetryableError(errors.New("foo")), false},
		{NotRetryableError(oerrors.New("foo")), false},
		{oerrors.Wrap(NotRetryableError(errors.New("foo")), "bar"), false},

		{httpError(http.StatusBadGateway), true},
		{oerrors.Wrap(httpError(http.StatusBadGateway), "bar"), true},
		{httpError(http.StatusInternalServerError), true},
		{oerrors.Wrap(httpError(http.StatusInternalServerError), "bar"), true},
		{httpError(http.StatusGatewayTimeout), true},
		{oerrors.Wrap(httpError(http.StatusGatewayTimeout), "bar"), true},
		{httpError(429), true},
		{oerrors.Wrap(httpError(429), "bar"), true},
		{httpError(http.StatusNotFound), false},
		{oerrors.Wrap(httpError(http.StatusNotFound), "bar"), false},
		{NotFoundError(errors.New("foo")), false},
		{httpError(http.StatusBadRequest), false},
		{oerrors.Wrap(httpError(http.StatusBadRequest), "bar"), false},
		{InvalidError(errors.New("foo")), false},
		{httpError(http.StatusForbidden), false},
		{oerrors.Wrap(httpError(http.StatusForbidden), "bar"), false},
		{ForbiddenError(errors.New("foo")), false},

		{NotRetryableError(httpError(http.StatusBadGateway)), false},
		{NotRetryableError(oerrors.Wrap(httpError(http.StatusBadGateway), "bar")), false},
		{NotRetryableError(httpError(http.StatusInternalServerError)), false},
		{NotRetryableError(oerrors.Wrap(httpError(http.StatusInternalServerError), "bar")), false},
		{NotRetryableError(httpError(http.StatusGatewayTimeout)), false},
		{NotRetryableError(oerrors.Wrap(httpError(http.StatusGatewayTimeout), "bar")), false},
		{NotRetryableError(httpError(429)), false},
		{NotRetryableError(oerrors.Wrap(httpError(429), "bar")), false},
		{NotRetryableError(httpError(http.StatusNotFound)), false},
		{NotRetryableError(oerrors.Wrap(httpError(http.StatusNotFound), "bar")), false},
		{NotRetryableError(NotFoundError(errors.New("foo"))), false},
		{NotRetryableError(httpError(http.StatusBadRequest)), false},
		{NotRetryableError(oerrors.Wrap(httpError(http.StatusBadRequest), "bar")), false},
		{NotRetryableError(InvalidError(errors.New("foo"))), false},
		{NotRetryableError(httpError(http.StatusForbidden)), false},
		{NotRetryableError(oerrors.Wrap(httpError(http.StatusForbidden), "bar")), false},
		{NotRetryableError(ForbiddenError(errors.New("foo"))), false},

		{retryable(false), false},
		{retryable(true), true},
		{NotRetryableError(retryable(false)), false},
		{NotRetryableError(retryable(true)), false},
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
	//err = NotRetryableError(err)
	IsRetryable(err) // will return false
}

func TestIsNotRetryable(t *testing.T) {
	tests := []struct {
		err  error
		want bool
	}{
		{nil, false},
		{RetryableError(nil), false},
		{NotRetryableError(nil), false},
		{oerrors.Wrap(RetryableError(nil), "bar"), false},
		{NotRetryableError(oerrors.Wrap(RetryableError(nil), "bar")), false},

		{errors.New("foo"), false},
		{oerrors.New("foo"), false},
		{oerrors.Wrap(errors.New("foo"), "bar"), false},

		{RetryableError(errors.New("foo")), false},
		{NotRetryableError(errors.New("foo")), true},
		{NotRetryableError(RetryableError(errors.New("foo"))), true},
		{RetryableError(oerrors.New("foo")), false},
		{NotRetryableError(RetryableError(oerrors.New("foo"))), true},
		{NotRetryableError(oerrors.New("foo")), true},
		{oerrors.Wrap(RetryableError(errors.New("foo")), "bar"), false},
		{oerrors.Wrap(NotRetryableError(errors.New("foo")), "bar"), true},
		{NotRetryableError(oerrors.Wrap(RetryableError(errors.New("foo")), "bar")), true},
		{RetryableError(NotRetryableError(oerrors.Wrap(RetryableError(errors.New("foo")), "bar"))), false},

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

func TestNewRetryableError(t *testing.T) {
	err := NewRetryableError("retryable error")
	if IsNotRetryable(err) {
		t.Errorf("NewRetryableError must return a retryable error")
	}
	if err.Error() != "retryable error" {
		t.Errorf("NewRetryableError: expected :%s, got %s", "test", err.Error())
	}
}

func ExampleNewRetryableError() {
	err := NewRetryableError("test")
	IsRetryable(err) // will return true
}

func TestNewRetryableErrorf(t *testing.T) {
	err := NewRetryableErrorf("retryable error %s", "formatted")
	if IsNotRetryable(err) {
		t.Errorf("NewRetryableErrorf must return a retryable error")
	}
	if err.Error() != "retryable error formatted" {
		t.Errorf("NewRetryableErrorf: expected :%s, got %s", "retryable error formatted", err.Error())
	}
}

func ExampleNewRetryableErrorf() {
	err := NewRetryableErrorf("Unable to read data for device %d", 70)
	IsRetryable(err) // will return true
}
