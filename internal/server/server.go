package server

import (
	"github.com/atyiadam/k8s-mutating-webhook/internal/handlers"
	"net/http"
)

type Server struct {
	*http.Server
	certFile string
	keyFile  string
}

func NewServer(addr string, certFile string, keyFile string) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/mutate/pod", handlers.MutatePod)

	return &Server{
		Server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		certFile: certFile,
		keyFile:  keyFile,
	}
}

func (s *Server) ListenAndServeTLS() error {
	return s.Server.ListenAndServeTLS(s.certFile, s.keyFile)
}
