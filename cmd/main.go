package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/abdulmanafc2001/url-shortner/pkg/api/server"
	"github.com/abdulmanafc2001/url-shortner/pkg/logger"
	"github.com/abdulmanafc2001/url-shortner/pkg/service"
	"github.com/abdulmanafc2001/url-shortner/pkg/storage"
)

func main() {
	var port string
	var baseURL string

	flag.StringVar(&port, "port", "8080", "The port to listen on")
	flag.StringVar(&baseURL, "baseurl", "http://localhost", "Base url for shortend url routes")

	// Initialize logger
	log := logger.NewLogger()
	store := storage.NewMemoryStore()
	shortenerService := service.NewShortenerService(store)

	// Create server config
	serverConfig := server.ResourceHandlersConfig{
		Logger:         log,
		ShortenService: shortenerService,
		BaseURL:        baseURL + ":" + port,
	}

	// Create and start server
	srv := server.NewServer(serverConfig)

	go func() {
		if err := srv.Start(port); err != nil {
			log.Error("Failed to start server", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info("Shutting down server...", nil)

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 15)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", map[string]interface{}{
			"error": err.Error(),
		})
		os.Exit(1)
	}

	log.Info("Server stopped gracefully", nil)
}
