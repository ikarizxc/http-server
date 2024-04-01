package server

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(adress, port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           adress + ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 28,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
