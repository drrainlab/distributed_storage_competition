package server

import (
	"context"
	"karma8/internal/server/handlers"
	"log"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

func NewServer(addr string, storageServiceHandler *handlers.StorageService) *Server {
	return &Server{
		server: &http.Server{
			Addr: addr,
			// Handler: handlers.CreateRouter(storageServiceHandler),
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		log.Printf("server started at %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Println("server stopped")
				return
			}
			log.Printf("Serving error: %s", err.Error())
		}
	}()

	select {
	case <-ctx.Done():
	}

	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	err := s.server.Shutdown(ctx)
	cancel()
	return err
}
