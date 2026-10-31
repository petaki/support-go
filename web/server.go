package web

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/petaki/support-go/cli"
)

// Server type.
type Server struct {
	*http.Server
	InfoLog *log.Logger
}

// NewServer function.
func NewServer(addr string, handler http.Handler, errorLog *log.Logger) *Server {
	return &Server{
		Server: &http.Server{
			Addr:         addr,
			ErrorLog:     errorLog,
			Handler:      handler,
			IdleTimeout:  time.Minute,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		InfoLog: cli.InfoLog,
	}
}

// Run function.
func (s *Server) Run(ctx context.Context) error {
	if ctx.Err() != nil {
		return nil
	}

	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	s.InfoLog.Printf("Starting server on "+cli.Green("%s"), listener.Addr())

	result := make(chan error, 1)

	go func() {
		result <- s.Server.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		s.InfoLog.Print("Server stopped")

		return nil
	case err := <-result:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return err
	}
}

// RunAndShutdown function.
func (s *Server) RunAndShutdown(ctx context.Context, timeout time.Duration) error {
	err := s.Run(ctx)
	if err != nil {
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = s.Shutdown(shutdownCtx)
	if err != nil {
		return err
	}

	s.InfoLog.Print("Server exited properly")

	return nil
}
