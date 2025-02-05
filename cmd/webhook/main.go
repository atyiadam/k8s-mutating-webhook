package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/atyiadam/k8s-mutating-webhook/internal/server"
	"github.com/atyiadam/k8s-mutating-webhook/pkg/utils"
)

func main() {
	certFile := os.Getenv("TLS_CERT_FILE")
	keyFile := os.Getenv("TLS_KEY_FILE")
	port := os.Getenv("PORT")

	if certFile == "" {
		certFile = "/home/adama/certs/webhook-server/webhook-server.crt"
	}
	if keyFile == "" {
		keyFile = "/home/adama/certs/webhook-server/webhook-server.key"
	}
	if port == "" {
		port = "8080"
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := server.NewServer(":"+port, certFile, keyFile)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	go func() {
		if err := s.ListenAndServeTLS(); err != http.ErrServerClosed {
			utils.LogError(err, "Server error")
			cancel()
		}
	}()

	<-ctx.Done()

	if err := s.Shutdown(ctx); err != nil {
		utils.LogError(err, "Shutdown error")
	}
}
