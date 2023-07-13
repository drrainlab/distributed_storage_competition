package server

import (
	"context"
	"errors"
	"karma8/internal/server/handlers"
	"log"
	"net/http"
)

type Server struct {
	server  *http.Server
	handler *handlers.Handler
}

func NewServer(addr string, h *handlers.Handler) *Server {
	return &Server{
		server: &http.Server{
			Addr:    addr,
			Handler: handlers.CreateRouter(h),
		},
		handler: h,
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.handler.Shutdown()
	return s.server.Shutdown(ctx)
}

func (s *Server) Run(ctx context.Context) {

	log.Println("running server at ", s.server.Addr)
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalln("serve listener err", err)
	}
}
