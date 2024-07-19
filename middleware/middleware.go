package middleware

import (
	"net/http"
)

// Middleware is a callback used to perform action before or after
// the handler is invoked. It returns a http.Handler to allow for middleware chaining.
//
// Middleware at some point must either invoke the given handler or
// respond to the request early in order to not break the execution flow.
//
// Middleware must not respond to the request after invoking the handler.
// It should assume that the controller has already been called and the
// response has been sent to the client.
type Middleware = func(http.Handler) http.Handler

// Apply applies middleware to the supplied handler and returns it.
// Middleware is applied in order it is provided, left to right.
//
// If the handler is a pointer it will be overwritten and should be
// discarded. The returned handler should be used instead.
func Apply(h http.Handler, middlewares ...Middleware) http.HandlerFunc {
	for _, middleware := range middlewares {
		h = middleware(h)
	}

	return h.ServeHTTP
}
