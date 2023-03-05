package main

import (
	"log"

	"github.com/tjper/gossip-glomers/internal/multibroadcast"
)

func main() {
	srv := multibroadcast.NewServer()
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
