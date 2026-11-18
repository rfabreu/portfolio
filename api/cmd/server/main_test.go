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
	// The write deadline starts once headers are read, so it covers the body read and the Gemini call.
	if srv.WriteTimeout <= srv.ReadTimeout+chat.RequestTimeout {
		t.Errorf("WriteTimeout (%v) must exceed ReadTimeout (%v) plus the Gemini timeout (%v)", srv.WriteTimeout, srv.ReadTimeout, chat.RequestTimeout)
	}
}
