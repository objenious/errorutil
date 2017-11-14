package backoffutil

import (
	"net/http"

	"github.com/objenious/errorutil"
)

func ExampleRetry() {
	Retry(func() error {
		resp, err := http.Get("http://www.example.com")
		if err != nil {
			return errorutil.RetryableError(err)
		}
		defer resp.Body.Close()
		if err := errorutil.HTTPError(resp); err != nil {
			return err
		}
		// Do something
		return nil
	})
}
