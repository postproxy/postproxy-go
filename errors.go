package postproxy

import "fmt"

// PostProxyError represents an error returned by the PostProxy API.
type PostProxyError struct {
	StatusCode int
	Message    string
	Response   map[string]any
}

func (e *PostProxyError) Error() string {
	if e.Response != nil {
		return fmt.Sprintf("PostProxy API error (status %d): %s — %v", e.StatusCode, e.Message, e.Response)
	}
	return fmt.Sprintf("PostProxy API error (status %d): %s", e.StatusCode, e.Message)
}

// IsAuthenticationError reports whether the error is a 401 authentication error.
func IsAuthenticationError(err error) bool {
	e, ok := err.(*PostProxyError)
	return ok && e.StatusCode == 401
}

// IsNotFoundError reports whether the error is a 404 not found error.
func IsNotFoundError(err error) bool {
	e, ok := err.(*PostProxyError)
	return ok && e.StatusCode == 404
}

// IsConflictError reports whether the error is a 409 conflict — a duplicate
// submission (Response["duplicate_post_id"]), a backfill that is already
// running (Response["profile_sync_id"]), or a request whose Idempotency-Key is
// still in flight.
func IsConflictError(err error) bool {
	e, ok := err.(*PostProxyError)
	return ok && e.StatusCode == 409
}

// IsValidationError reports whether the error is a 422 validation error.
func IsValidationError(err error) bool {
	e, ok := err.(*PostProxyError)
	return ok && e.StatusCode == 422
}

// IsBadRequestError reports whether the error is a 400 bad request error.
func IsBadRequestError(err error) bool {
	e, ok := err.(*PostProxyError)
	return ok && e.StatusCode == 400
}
