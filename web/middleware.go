package web

import (
	"fmt"
	"net/http"
)

// RecoverPanic function.
func RecoverPanic(onError func(w http.ResponseWriter, r *http.Request, err error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				recovered := recover()
				if recovered == nil {
					return
				}

				if recovered == http.ErrAbortHandler {
					panic(recovered)
				}

				w.Header().Set("Connection", "close")

				if err, ok := recovered.(error); ok {
					onError(w, r, fmt.Errorf("%w: %w", ErrPanic, err))

					return
				}

				onError(w, r, fmt.Errorf("%w: %v", ErrPanic, recovered))
			}()

			next.ServeHTTP(w, r)
		})
	}
}
