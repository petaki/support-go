package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoverPanic(t *testing.T) {
	t.Run("Panic", func(t *testing.T) {
		var got error

		handler := RecoverPanic(func(w http.ResponseWriter, r *http.Request, err error) {
			got = err
			w.WriteHeader(http.StatusInternalServerError)
		})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("boom")
		}))

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))

		if !errors.Is(got, ErrPanic) {
			t.Errorf("expected: %v, got: %v", ErrPanic, got)
		}

		if rr.Header().Get("Connection") != "close" {
			t.Errorf("expected Connection: close, got: %v", rr.Header().Get("Connection"))
		}
	})

	t.Run("Error Panic", func(t *testing.T) {
		cause := errors.New("cause")

		var got error

		handler := RecoverPanic(func(w http.ResponseWriter, r *http.Request, err error) {
			got = err
		})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(cause)
		}))

		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

		if !errors.Is(got, ErrPanic) || !errors.Is(got, cause) {
			t.Errorf("expected: %v and %v, got: %v", ErrPanic, cause, got)
		}
	})

	t.Run("Abort Handler", func(t *testing.T) {
		handler := RecoverPanic(func(w http.ResponseWriter, r *http.Request, err error) {
			t.Error("expected onError not to be called")
		})(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic(http.ErrAbortHandler)
		}))

		defer func() {
			if recovered := recover(); recovered != http.ErrAbortHandler {
				t.Errorf("expected: %v, got: %v", http.ErrAbortHandler, recovered)
			}
		}()

		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
	})
}
