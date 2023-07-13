package main

import (
	"context"
	"fmt"
	"io"
	"karma8/internal/service"
	"log"
	"strings"
)

func main() {

	cfg := service.Config{
		NodesNum: 6,
		Capacity: 10000,
	}

	objectStorage, err := service.NewService(cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println(objectStorage.Nodes())

	strR := strings.NewReader(`When doneCh is closed above, the function asyncReader will return. The goroutine we created will also return the next time it evaluates the condition in the for loop. But, what if the goroutine is blocking on r.Read()? Then, we essentially have leaked a goroutine. We’re stuck until the reader unblocks.`)

	ctx := context.Background()
	if err := objectStorage.Store(ctx, "test", uint64(strR.Len()), strR); err != nil {
		log.Fatalf("error storing object: %v\n", err)
	}

	fileReader, err := objectStorage.Load(ctx, "test")
	if err != nil {
		log.Fatalf("error loading object: %v\n", err)
	}

	// Read data from multi reader
	b, err := io.ReadAll(fileReader)

	if err != nil {
		panic(err)
	}

	// Optional: Verify data
	fmt.Println(string(b))

}
