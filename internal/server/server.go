package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	srv             *http.Server
	shutdownTimeout time.Duration
}

func NewServer(port string, handler http.Handler, readTimeout, writeTimeout, shutdownTimeout time.Duration) *Server {
	return &Server{
		srv: &http.Server{
			Addr:         ":" + port,
			Handler:      handler,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  120 * time.Second,
		},
		shutdownTimeout: shutdownTimeout,
	}
}

func (s *Server) RunServer() error {
	slog.Info("starting HTTP server", "addr", s.srv.Addr)
	err := s.srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (s *Server) Shutdown() error {
	slog.Info("starting graceful shutdown of http server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		if closeErr := s.srv.Close(); closeErr != nil {
			return fmt.Errorf("failed to force-shutdown HTTP server: %w", closeErr)
		}
		return fmt.Errorf("failed to shutdown HTTP server gracefully: %w", err)
	}

	return nil
}
