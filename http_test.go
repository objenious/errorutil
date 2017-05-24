package errorutil

import (
	"errors"
	"net/http"
)

func ExampleHTTPError() {
	resp, err := http.Get("http://www.example.com")
	if err != nil {
		// return error
	}
	defer resp.Body.Close()
	if err := HTTPError(resp); err != nil {
		// return error
	}
	// handle response
}

func ExampleNotFoundError() {
	var w http.ResponseWriter
	err := errors.New("some error")
	err = NotFoundError(err)
	w.WriteHeader(HTTPStatusCode(err)) // returns http.StatusNotFound
}

func ExampleForbiddenError() {
	var w http.ResponseWriter
	err := errors.New("some error")
	err = ForbiddenError(err)
	w.WriteHeader(HTTPStatusCode(err)) // returns http.StatusForbidden
}

func ExampleInvalidError() {
	var w http.ResponseWriter
	err := errors.New("some error")
	err = InvalidError(err)
	w.WriteHeader(HTTPStatusCode(err)) // returns http.StatusBadRequest
}
