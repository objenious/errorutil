package errorutil

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"testing"

	pkgerrors "github.com/pkg/errors"
)

func TestHTTPStatusCode(t *testing.T) {
	tests := []struct {
		err  error
		want int
	}{
		{nil, http.StatusOK},

		{errors.New("foo"), http.StatusInternalServerError},

		{os.ErrNotExist, http.StatusNotFound},
		{pkgerrors.Wrap(os.ErrNotExist, "bar"), http.StatusNotFound},
		{os.ErrPermission, http.StatusForbidden},
		{pkgerrors.Wrap(os.ErrPermission, "bar"), http.StatusForbidden},
		{sql.ErrNoRows, http.StatusNotFound},
		{pkgerrors.Wrap(sql.ErrNoRows, "bar"), http.StatusNotFound},

		{httpError(http.StatusNotFound), http.StatusNotFound},
		{pkgerrors.Wrap(httpError(http.StatusNotFound), "bar"), http.StatusNotFound},
		{httpError(http.StatusForbidden), http.StatusForbidden},
		{pkgerrors.Wrap(httpError(http.StatusForbidden), "bar"), http.StatusForbidden},

		{NotFoundError(errors.New("foo")), http.StatusNotFound},
		{pkgerrors.Wrap(NotFoundError(errors.New("foo")), "bar"), http.StatusNotFound},
		{ForbiddenError(errors.New("foo")), http.StatusForbidden},
		{pkgerrors.Wrap(ForbiddenError(errors.New("foo")), "bar"), http.StatusForbidden},
		{InvalidError(errors.New("foo")), http.StatusBadRequest},
		{pkgerrors.Wrap(InvalidError(errors.New("foo")), "bar"), http.StatusBadRequest},
		{ConflictError(errors.New("foo")), http.StatusConflict},
		{pkgerrors.Wrap(ConflictError(errors.New("foo")), "bar"), http.StatusConflict},
	}
	for _, tt := range tests {
		got := HTTPStatusCode(tt.err)
		if got != tt.want {
			t.Errorf("HTTPStatusCode(%q): got: %v, want %v", tt.err, got, tt.want)
		}
	}
}

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

func ExampleConflictError() {
	var w http.ResponseWriter
	err := errors.New("some error")
	err = ConflictError(err)
	w.WriteHeader(HTTPStatusCode(err)) // returns http.StatusConflict
}
