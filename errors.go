package envoyerapi

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrBadRequest         = errors.New("api: bad request")
	ErrUnauthenticated    = errors.New("api: unauthenticated")
	ErrForbidden          = errors.New("api: forbidden")
	ErrNotFound           = errors.New("api: not found")
	ErrUnprocessable      = errors.New("api: payload missing or invalid")
	ErrTooManyRequests    = errors.New("api: too many attempts")
	ErrServer             = errors.New("api: server error")
	ErrServiceUnavailable = errors.New("api: service unavailable")
)

type HTTPError struct {
	Status      int
	Method      string
	URL         string
	BodySnippet string
	Headers     http.Header
	Cause       error
}

func (e *HTTPError) Error() string {
	if e.BodySnippet != "" {
		return fmt.Sprintf("%s %s: %d (%s)", e.Method, e.URL, e.Status, e.BodySnippet)
	}

	return fmt.Sprintf("%s %s: %d", e.Method, e.URL, e.Status)
}

func (e *HTTPError) Unwrap() error {
	return e.Cause
}

func (e *HTTPError) Is(target error) bool {
	switch {
	case errors.Is(target, ErrBadRequest):
		return e.Status == http.StatusBadRequest
	case errors.Is(target, ErrUnauthenticated):
		return e.Status == http.StatusUnauthorized
	case errors.Is(target, ErrForbidden):
		return e.Status == http.StatusForbidden
	case errors.Is(target, ErrNotFound):
		return e.Status == http.StatusNotFound
	case errors.Is(target, ErrUnprocessable):
		return e.Status == http.StatusUnprocessableEntity
	case errors.Is(target, ErrTooManyRequests):
		return e.Status == http.StatusTooManyRequests
	case errors.Is(target, ErrServiceUnavailable):
		return e.Status == http.StatusServiceUnavailable
	case errors.Is(target, ErrServer):
		return e.Status >= 500 && e.Status <= 599
	default:
		return false
	}
}

func NewHTTPError(req *http.Request, resp *http.Response, err *error) *HTTPError {
	defer io.Copy(io.Discard, resp.Body)

	var body []byte
	if resp.Body != nil {
		var readErr error
		body, readErr = io.ReadAll(io.LimitReader(resp.Body, 8<<10)) // 8 KiB snippet
		if readErr != nil {
			body = []byte("[error capturing body: " + readErr.Error() + "]")
		}
	}

	return &HTTPError{
		Status:      resp.StatusCode,
		Method:      req.Method,
		URL:         req.URL.String(),
		BodySnippet: string(body),
		Headers:     resp.Header.Clone(),
		Cause:       *err,
	}
}
