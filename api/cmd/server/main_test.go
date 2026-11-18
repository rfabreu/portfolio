package main

import (
	"net/http"
	"testing"

	"github.com/rfabreu/portfolio-api/internal/chat"
)

func TestNewServer_SetsTimeouts(t *testing.T) {
	srv := newServer(":0", http.NotFoundHandler())

	if srv.ReadHeaderTimeout <= 0 || srv.ReadTimeout <= 0 || srv.WriteTimeout <= 0 || srv.IdleTimeout <= 0 {
		t.Errorf("all server timeouts must be set, got %+v", srv)
	}
	if srv.WriteTimeout <= chat.RequestTimeout {
		t.Errorf("WriteTimeout (%v) must exceed the Gemini timeout (%v)", srv.WriteTimeout, chat.RequestTimeout)
	}
}
