package errorutil

import "net/http"

// HTTPStatusCodeEr defines errors that should return a specific HTTP status code
type HTTPStatusCodeEr interface {
	HTTPStatusCode() int
}

// HTTPStatusCode returns the status code that an error should return.
// If the error is nil, StatusOK is returned.
// If the error does not implement HTTPStatusCodeEr, StatusInternalServerError is returned.
func HTTPStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if status, ok := err.(HTTPStatusCodeEr); ok {
		return status.HTTPStatusCode()
	}
	return http.StatusInternalServerError
}

// HTTPError builds an error based on a http.Response. If status code is < 300, nil is returned.
// Otherwise, errors implementing the various interfaces (Retryabler, HTTPStatusCodeEr) are returned
func HTTPError(resp *http.Response) error {
	if resp.StatusCode < 300 {
		return nil
	}
	return httpError(resp.StatusCode)
}

type httpError int

func (err httpError) Error() string {
	switch err {
	case 429:
		return "Too Many Requests"
	default:
		return http.StatusText(int(err))
	}
}

func (err httpError) HTTPStatusCode() int {
	return int(err)
}

func (err httpError) Retryable() bool {
	switch int(err) {
	case http.StatusBadGateway, http.StatusGatewayTimeout, http.StatusServiceUnavailable, http.StatusInternalServerError:
		return true
	case 429:
		return true
	default:
		return false
	}
}

// NotFoundError marks an error as "not found". The calling http handler
// should return a StatusNotFound status code. It returns nil if the error is nil.
func NotFoundError(err error) error {
	if err == nil {
		return nil
	}
	return &notFoundError{err: err}
}

type notFoundError struct {
	err error
}

func (err *notFoundError) Error() string {
	return err.err.Error()
}

func (err *notFoundError) HTTPStatusCode() int {
	return http.StatusNotFound
}

func (err *notFoundError) IsRetryable() bool {
	return false
}

func (err *notFoundError) Cause() error {
	return err.err
}

// ForbiddenError marks an error as "access forbidden". The calling http handler
// should return a StatusForbidden status code. It returns nil if the error is nil.
func ForbiddenError(err error) error {
	if err == nil {
		return nil
	}
	return &forbiddenError{err: err}
}

type forbiddenError struct {
	err error
}

func (err *forbiddenError) Error() string {
	return err.err.Error()
}

func (err *forbiddenError) HTTPStatusCode() int {
	return http.StatusForbidden
}

func (err *forbiddenError) IsRetryable() bool {
	return false
}

func (err *forbiddenError) Cause() error {
	return err.err
}

// InvalidError marks an error as "invalid". The calling http handler
// should return a StatusBadRequest status code. It returns nil if the error is nil.
func InvalidError(err error) error {
	if err == nil {
		return nil
	}
	return &invalidError{err: err}
}

type invalidError struct {
	err error
}

func (err *invalidError) Error() string {
	return err.err.Error()
}

func (err *invalidError) HTTPStatusCode() int {
	return http.StatusBadRequest
}

func (err *invalidError) IsRetryable() bool {
	return false
}

func (err *invalidError) Cause() error {
	return err.err
}
