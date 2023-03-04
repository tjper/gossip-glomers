package main

import (
	"log"

	"github.com/tjper/gossip-glomers/internal/idgen"
)

func main() {
	srv := idgen.NewServer()
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
