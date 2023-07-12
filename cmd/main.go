package main

import (
	"context"
	"fmt"
	"karma8/internal/service"
	"log"
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

	if err := objectStorage.Store(context.Background(), "test", 6, nil); err != nil {
		log.Printf("error storing object: %v\n", err)
	}

}
