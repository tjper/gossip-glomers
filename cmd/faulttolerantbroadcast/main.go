package main

import (
	"log"

	"github.com/tjper/gossip-glomers/internal/faulttolerantbroadcast"
)

func main() {
	srv := faulttolerantbroadcast.NewServer()
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
