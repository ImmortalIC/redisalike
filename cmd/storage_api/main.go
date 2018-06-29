package main

import (
	"context"
	"fmt"
	"github.com/ImmortalIC/redisalike/api"
	"github.com/ImmortalIC/redisalike/storage"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := new(sync.WaitGroup)
	storage.Init(ctx, wg)
	err := api.StartServer()
	fmt.Println(err)
	cancel()
	wg.Wait()
}
