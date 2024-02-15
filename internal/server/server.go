package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andreevym/gofermart/pkg/logger"
	"go.uber.org/zap"
)

type Server struct {
	Handler http.Handler
	Server  *http.Server
}

func NewServer(handler http.Handler) *Server {
	return &Server{
		Handler: handler,
	}
}

func (s *Server) Run(addr string) {
	s.Server = &http.Server{Addr: addr, Handler: s.Handler}

	logger.Logger().Info("Server listening", zap.String("addr", addr))
	if err := s.Server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func (s *Server) Shutdown() {
	if s.Server == nil {
		panic("server can't be nil")
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	fmt.Println("Server stopped gracefully")
}
