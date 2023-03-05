package main

import (
	"log"

	"github.com/tjper/gossip-glomers/internal/broadcast"
)

func main() {
	srv := broadcast.NewServer()
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
