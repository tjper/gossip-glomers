package main

import (
	"log"

	"github.com/tjper/gossip-glomers/internal/echo"
)

func main() {
	srv := echo.NewServer()
	if err := srv.Run(); err != nil {
		log.Fatalf("while running echo server: %s", err)
	}
}
