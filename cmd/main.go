package main

import (
	"context"
	"karma8/internal/server"
	"karma8/internal/server/handlers"
	"karma8/internal/service"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	ctx := context.Background()

	cfg := service.Config{
		NodesNum: 6,
		Capacity: 10000,
	}

	objectStorage, err := service.NewService(cfg)
	if err != nil {
		panic(err)
	}

	srv := server.NewServer(":8080", handlers.NewHandler(objectStorage))

	go srv.Run(ctx)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	srv.Shutdown(ctx)

}
