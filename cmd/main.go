package main

import (
	"context"
	"fmt"
	"karma8/internal/service"
	"log"
	"strings"
	"time"
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

	strR := strings.NewReader("aabbccddeefff")

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)
	defer cancel()
	if err := objectStorage.Store(ctx, "test", uint64(strR.Len()), strR); err != nil {
		log.Printf("error storing object: %v\n", err)
	}

}
