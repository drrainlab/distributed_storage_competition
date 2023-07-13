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

	strR := strings.NewReader("111111222222333333444444555555666666")

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
