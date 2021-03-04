package main

import (
	"log"

	"github.com/hiqua/rworker/internal"
)

func main() {
	log.Fatal(server.Serve())
}
