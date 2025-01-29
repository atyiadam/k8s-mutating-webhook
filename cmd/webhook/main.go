package main

import (
	"log"
	"os"

	"github.com/atyiadam/k8s-mutating-webhook/internal/server"
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

	s := server.NewServer(":"+port, certFile, keyFile)

	log.Fatal(s.ListenAndServeTLS())
}
