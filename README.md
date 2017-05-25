# errorutil [![Travis-CI](https://travis-ci.org/objenious/errorutil.svg)](https://travis-ci.org/objenious/errorutil)  [![GoDoc](https://godoc.org/github.com/objenious/errorutil?status.svg)](http://godoc.org/github.com/objenious/errorutil)

`go get github.com/objenious/errorutil`

## Retryable errors

```go
func foo() error {
  err := bar()
  if err != nil {
    // it should be retried
    return errorutil.RetryableError(err)
  }
  return baz() // will not be retried
}

func main() {
  err := foo()
  if errorutil.IsRetryable(err) {
    // retry !
  }
}
```

## HTTP Aware errors

Build an error based on a http.Response. Status code above 299, except 304, will be considered an error.

It will be retryable if status code is http.StatusBadGateway, http.StatusGatewayTimeout, http.StatusServiceUnavailable, http.StatusInternalServerError or 429 (Too many requests).

```go
resp, err := http.Get("http://www.example.com")
if err != nil {
  // return error
}
defer resp.Body.Close()
if err := errorutil.HTTPError(resp); err != nil {
  // return error
}
// handle response
```

Find the most appropriate status code for an error :

```go
w.WriteHeader(errorutil.HTTPStatusCode(err))
```

Generate specific error types :

```go
err := errors.New("some error")
err = errorutil.NotFoundError(err)
w.WriteHeader(errorutil.HTTPStatusCode(err)) // returns http.StatusNotFound
```
## Exponential backoff

```go
backoffutil.Retry(func() error {
  resp, err := http.Get("http://www.example.com")
  if err != nil {
    return err
  }
  defer resp.Body.Close()
  if err := errorutil.HTTPError(resp); err != nil {
    return err
  }
  // Do something
  return nil
})
```

## Notes

errorutil is compatible with https://github.com/pkg/errors :

```go
err = errors.Wrap(errorutil.RetraybleError(err), "some message")
errorutil.IsRetryable(err) // returns true
```
