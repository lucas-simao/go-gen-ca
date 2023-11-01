package main

import (
	"net/http"

	"github.com/lucas-simao/go-gen-ca/internal/server"
	goLog "github.com/lucas-simao/golog"
)

func main() {
	if err := http.ListenAndServe(":3000", server.NewServer()); err != nil {
		goLog.Critical(err)
	}
}
