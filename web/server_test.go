package web

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"testing"
	"time"
)

func newTestServer(addr string) *Server {
	server := NewServer(addr, http.NotFoundHandler(), log.New(io.Discard, "", 0))
	server.InfoLog = log.New(io.Discard, "", 0)

	return server
}

func waitRun(t *testing.T, run func() error) error {
	t.Helper()

	result := make(chan error, 1)

	go func() {
		result <- run()
	}()

	select {
	case err := <-result:
		return err
	case <-time.After(time.Second):
		t.Fatal("expected run to return within 1s")
	}

	return nil
}

func TestNewServer(t *testing.T) {
	server := newTestServer(":4000")

	if server.Addr != ":4000" {
		t.Errorf("expected: :4000, got: %v", server.Addr)
	}

	if server.IdleTimeout != time.Minute || server.ReadTimeout != 5*time.Second || server.WriteTimeout != 10*time.Second {
		t.Errorf("unexpected timeouts: %v, %v, %v", server.IdleTimeout, server.ReadTimeout, server.WriteTimeout)
	}
}

func TestServer_Run(t *testing.T) {
	server := newTestServer("127.0.0.1:0")
	ctx, cancel := context.WithCancel(context.Background())

	time.AfterFunc(50*time.Millisecond, cancel)

	err := waitRun(t, func() error {
		return server.Run(ctx)
	})
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	err = server.Shutdown(context.Background())
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestServer_RunCancelledContext(t *testing.T) {
	server := newTestServer("127.0.0.1:0")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := server.Run(ctx)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestServer_RunListenError(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	server := newTestServer(listener.Addr().String())

	err = waitRun(t, func() error {
		return server.Run(context.Background())
	})
	if err == nil {
		t.Error("expected an error")
	}
}

func TestServer_RunAndShutdown(t *testing.T) {
	server := newTestServer("127.0.0.1:0")
	ctx, cancel := context.WithCancel(context.Background())

	time.AfterFunc(50*time.Millisecond, cancel)

	err := waitRun(t, func() error {
		return server.RunAndShutdown(ctx, time.Second)
	})
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}
