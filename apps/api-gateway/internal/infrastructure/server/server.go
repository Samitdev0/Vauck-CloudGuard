package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type Server struct {
	httpServer *http.Server
	logger     zerolog.Logger
}

func New(
	host string,
	port int,
	handler http.Handler,
	logger zerolog.Logger,
) *Server {

	address := fmt.Sprintf("%s:%d", host, port)

	return &Server{
		httpServer: &http.Server{
			Addr:              address,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
		logger: logger,
	}
}

func (s *Server) Start() error {

	s.logger.Info().
		Str("address", s.httpServer.Addr).
		Msg("starting api gateway")

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}

func (s *Server) HTTPServer() *http.Server {
	return s.httpServer
}
