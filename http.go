package errorutil

import "net/http"

// HTTPStatusCodeEr defines errors that should return a specific HTTP status code
type HTTPStatusCodeEr interface {
	HTTPStatusCode() int
}

// HTTPStatusCode returns the status code that a HTTP handler should return.
//
// If the error is nil, StatusOK is returned.
//
// If the error implements HTTPStatusCodeEr, it returns the corresponding status code.
//
// It tries to check some stdlib errors (testing the error string, to avoid importing unwanted packages),
// and returns appropriate status codes.
//
// Otherwise, StatusInternalServerError is returned.
func HTTPStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	type causer interface {
		Cause() error
	}

	for err != nil {
		if status, ok := err.(HTTPStatusCodeEr); ok {
			return status.HTTPStatusCode()
		}
		// Check errors from stdlib. Test string to avoid importing packages
		switch err.Error() {
		// package os
		case "permission denied":
			return http.StatusForbidden
		case "file does not exist":
			return http.StatusNotFound
		case "storage: bucket doesn't exist":
			return http.StatusNotFound
		case "storage: object doesn't exist":
			return http.StatusNotFound
		// package database/sql
		case "sql: no rows in result set":
			return http.StatusNotFound
		}
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return http.StatusInternalServerError
}

// HTTPError builds an error based on a http.Response. If status code is < 300 or 304, nil is returned.
// Otherwise, errors implementing the various interfaces (Retryabler, HTTPStatusCodeEr) are returned
func HTTPError(resp *http.Response) error {
	if resp.StatusCode < 300 || resp.StatusCode == http.StatusNotModified {
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

func (err httpError) StatusCode() int {
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

// ConflictError marks an error as "conflict". The calling http handler
// should return a Conflict status code. It returns nil if the error is nil.
func ConflictError(err error) error {
	if err == nil {
		return nil
	}
	return &conflictError{err: err}
}

type conflictError struct {
	err error
}

func (err *conflictError) Error() string {
	return err.err.Error()
}

func (err *conflictError) HTTPStatusCode() int {
	return http.StatusConflict
}

func (err *conflictError) IsRetryable() bool {
	return false
}

func (err *conflictError) Cause() error {
	return err.err
}
