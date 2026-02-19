package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/maklean/flux/server/api"
	"github.com/maklean/flux/server/server_interface"
)

func init() {
	// load env variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load env: %s", err.Error())
	}
}

func main() {
	go api.StartAPIServer()
	go server_interface.StartgRPCServer()

	select {}
}
