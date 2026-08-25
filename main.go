package main

import (
	"cairn-sonar/internal/handler"
	"cairn-sonar/internal/service"
	"cairn-sonar/internal/store"
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	addr := flag.String("addr", ":8097", "HTTP listen address")
	dbPath := flag.String("db", "cairn-sonar.db", "SQLite database path")
	flag.Parse()
	repo, err := store.Open(filepath.Clean(*dbPath))
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()
	svc := service.New(repo, service.Config{})
	defer svc.Close()
	logger := log.New(os.Stdout, "cairn-sonar ", log.LstdFlags|log.Lmicroseconds)
	srv := &http.Server{Addr: *addr, Handler: handler.New(svc, logger).Handler()}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		logger.Printf("listening on %s", *addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Printf("server: %v", err)
		}
	}()
	<-ctx.Done()
	_ = srv.Shutdown(context.Background())
}
