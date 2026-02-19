package main

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
	"github.com/maklean/flux/server/api"
	"github.com/maklean/flux/server/server_interface"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load env: %s", err.Error())
	}
}

func main() {
	var wg sync.WaitGroup

	wg.Go(api.StartAPIServer)
	wg.Go(server_interface.StartgRPCServer)

	wg.Wait()
}
