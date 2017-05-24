// Package errorutil allows errors to be tagged, allowing calling code to make decisions,
// such as "is the error retryable ?" or "what kind of HTTP status code should I return ?".
//
// It is compatible with github.com/pkg/errors :
//   err = errorutil.Retryable(err)
//   err = errors.Wrap(err, "an error occured")
//   errorutil.IsRetryable(err) // returns true
package errorutil
