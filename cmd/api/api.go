package main

import (
	"github.com/artrey/remux/pkg/middleware/logger"
	"github.com/artrey/remux/pkg/middleware/recoverer"
	"github.com/artrey/simple-pages/pkg/service"
	"github.com/artrey/simple-pages/pkg/storage/memory"
	"log"
	"net"
	"net/http"
	"os"
)

const (
	defaultHost = "0.0.0.0"
	defaultPort = "9999"
)

func main() {
	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = defaultHost
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultPort
	}

	address := net.JoinHostPort(host, port)
	log.Println(address)

	if err := execute(address); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func execute(address string) error {
	storage := memory.New()

	application := service.New(storage)
	application.Init(logger.Logger, recoverer.Recoverer)

	server := http.Server{
		Addr:    address,
		Handler: application,
	}
	return server.ListenAndServe()
}
