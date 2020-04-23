package main

import (
	"log"

	"github.com/psucodervn/go/logger"
	"github.com/psucodervn/go/server"
)

func main() {
	logger.Init(true, true)

	e := server.NewDefaultEchoServer()
	log.Fatal(e.Start(":8080"))
}
