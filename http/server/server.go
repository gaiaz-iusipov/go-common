package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

func New(addr string, handler http.Handler) Server {
	return Server{
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
		},
	}
}

type Server struct {
	httpServer *http.Server
}

func (s Server) Run(ctx context.Context) error {
	listener, err := new(net.ListenConfig).Listen(ctx, "tcp", s.httpServer.Addr)
	if err != nil {
		return fmt.Errorf("net listen: %w", err)
	}

	go func() {
		if serveErr := s.httpServer.Serve(listener); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			slog.ErrorContext(ctx, "failed to serve http server", "error", serveErr)
		}
	}()

	return nil
}

func (s Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
