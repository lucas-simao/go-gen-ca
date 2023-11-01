package main

import (
	"log"
	"net/http"

	"github.com/lucas-simao/go-gen/internal/server"
)

func main() {
	if err := http.ListenAndServe(":8080", server.NewServer()); err != nil {
		log.Panic(err)
	}
}
