package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	srv *http.Server
}

func NewServer(port string, handler http.Handler) *Server {
	return &Server{
		srv: &http.Server{
			Addr:    ":" + port,
			Handler: handler,
		},
	}
}

func (s *Server) RunServer(ctx context.Context) <-chan error {
	slog.Info("Starting server", "addr", s.srv.Addr)
	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)

		go func() {
			if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errChan <- err
			}
		}()

		<-ctx.Done()
		slog.InfoContext(ctx, "received a signal; draining HTTP server connection")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.srv.Shutdown(shutdownCtx); err != nil {
			slog.ErrorContext(shutdownCtx, "failed to shutdown HTTP server gracefully", "err", err)
			if closeErr := s.srv.Close(); closeErr != nil {
				slog.ErrorContext(shutdownCtx, "failed to force-shutdown HTTP server", "err", closeErr)
			}
		}
	}()

	return errChan
}
