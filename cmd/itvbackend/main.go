package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ahamtat/itvbackend/internal/app/storage/database"
	"github.com/ahamtat/itvbackend/internal/app/storage/memory"

	"github.com/ahamtat/itvbackend/internal/app/fetcher"

	"github.com/ahamtat/itvbackend/internal/app/server"
	"github.com/sirupsen/logrus"
)

var (
	port     string
	dsn      string
	mode     string
	timeout  int
	poolSize int
	logger   = logrus.New()
)

func init() {
	flag.StringVar(&port, "port", "8080", "server port number")
	flag.StringVar(&dsn, "dsn", "postgres://postgres:postgres@localhost:5432/itvbackend?sslmode=disable", "database connection string")
	flag.StringVar(&mode, "mode", "memory", "storage mode [memory, database]")
	flag.IntVar(&timeout, "timeout", 5, "timeout for external resource")
	flag.IntVar(&poolSize, "pool", 5, "size of worker pool & database connection pool")
	flag.Parse()
}

func main() {
	// Create application main context
	ctx, cancel := context.WithCancel(context.Background())

	var handler http.Handler
	switch mode {
	case "memory":
		handler = server.NewServer(
			fetcher.NewHTTPFetcher(time.Duration(timeout)*time.Second),
			memory.NewMemoryStorage())
	case "database":
		db, err := database.CreateDatabase(dsn, poolSize)
		defer func() { _ = db.Close() }()

		if err != nil {
			logger.Fatalf("failed creating database connection: %v\n", err)
		}
		handler = server.NewConcurrentServer(
			poolSize,
			fetcher.NewHTTPFetcher(time.Duration(timeout)*time.Second),
			database.NewDatabaseStorage(ctx, db))
	default:
		logger.Fatalf("wrong storage mode: %s\n", mode)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}

	// Start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("error listening HTTP server: %s", err.Error())
		}
	}()

	// Set interrupt handler
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Infoln("Application started. Press Ctrl+C to exit...")

	// Wait until user interrupt
	<-done

	// Cancel main context
	cancel()

	// Make server graceful shutdown
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if concurSrv, ok := handler.(*server.ConcurrentServer); ok {
		concurSrv.Close()
	}
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server shutdown failed: %v", err)
	}

	logger.Info("Application exited properly")
}
