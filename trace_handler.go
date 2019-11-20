package observability

import (
	"net/http"
)

// traceHeader adds a 'trace' header to response associated with trace ID of
// request.
func traceHeader(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := traceID(r.Context())
		w.Header().Set("trace", id)
		h.ServeHTTP(w, r)
	})
}
