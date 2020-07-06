package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ahamtat/itvbackend/internal/app/fetcher"

	"github.com/ahamtat/itvbackend/internal/app/server"
	"github.com/ahamtat/itvbackend/internal/app/storage"
	"github.com/sirupsen/logrus"
)

var (
	port    string
	timeout int
)

func main() {
	flag.StringVar(&port, "port", "8080", "server port number")
	flag.IntVar(&timeout, "timeout", 5, "timeout for external resource")
	flag.Parse()

	logger := logrus.New()

	srv := server.NewServer(
		fetcher.NewHTTPFetcher(time.Duration(timeout)*time.Second),
		storage.NewMemoryStorage(),
		logrus.New())

	// Start server
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%s", port), srv); err != nil {
			logger.Fatalf("error listening HTTP server: %s", err.Error())
		}
	}()

	// Set interrupt handler
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Infoln("Application started. Press Ctrl+C to exit...")

	// Wait until user interrupt
	<-done

	logger.Info("Application exited properly")
}
