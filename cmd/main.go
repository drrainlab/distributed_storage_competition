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

	// fmt.Println(objectStorage.Nodes())

	// strR := strings.NewReader(`When doneCh is closed above, the function asyncReader will return. The goroutine we created will also return the next time it evaluates the condition in the for loop. But, what if the goroutine is blocking on r.Read()? Then, we essentially have leaked a goroutine. We’re stuck until the reader unblocks.`)

	// ctx := context.Background()
	// if err := objectStorage.Store(ctx, "test", uint64(strR.Len()), strR); err != nil {
	// 	log.Fatalf("error storing object: %v\n", err)
	// }

	// fileReader, err := objectStorage.Load(ctx, "test")
	// if err != nil {
	// 	log.Fatalf("error loading object: %v\n", err)
	// }

	// // Read data from multi reader
	// b, err := io.ReadAll(fileReader)

	// if err != nil {
	// 	panic(err)
	// }

	// // Optional: Verify data
	// fmt.Println(string(b))

}
