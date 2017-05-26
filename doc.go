/*
Package errorutil allows errors to be tagged, allowing calling code to make decisions,
such as "is the error retryable ?" or "what kind of HTTP status code should I return ?".

Retryable errors

Mark errors as retryable :

  func foo() error {
    err := bar()
    if err != nil {
      // should be retried
      return errorutil.RetryableError(err)
    }
    return baz() // should not be retried
  }

  func main() {
    err := foo()
    if errorutil.IsRetryable(err) {
      // retry !
    }
  }

HTTP Aware errors

Build an error based on a http.Response. It will be retryable of status code is http.StatusBadGateway, http.StatusGatewayTimeout, http.StatusServiceUnavailable, http.StatusInternalServerError or 429 (Too many requests).

  resp, err := http.Get("http://www.example.com")
  if err != nil {
    // return error
  }
  defer resp.Body.Close()
  if err := errorutil.HTTPError(resp); err != nil {
    // return error
  }
  // handle response
  Find the most appropriate status code for an error :

  w.WriteHeader(errorutil.HTTPStatusCode(err))

Generate specific error types :

  err := errors.New("some error")
  err = errorutil.NotFoundError(err)
  w.WriteHeader(errorutil.HTTPStatusCode(err)) // returns http.StatusNotFound

Exponential backoff

see backoffutil sub package

Notes

errorutil is compatible with https://github.com/pkg/errors :

  err = errors.Wrap(errorutil.RetraybleError(err), "some message")
  errorutil.IsRetryable(err) // returns true
*/
package errorutil
