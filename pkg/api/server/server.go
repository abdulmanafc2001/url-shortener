package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/abdulmanafc2001/url-shortner/pkg/api/handlers"
	"github.com/abdulmanafc2001/url-shortner/pkg/logger"
)

type Server struct {
	httpServer *http.Server
	router     *Router
	logger     *logger.Logger
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace * with specific domain in production
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Refresh-Token, Tenant, Client-Secret, Client-ID, X-Config")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type ResourceHandlersConfig struct {
	Logger *logger.Logger
}

func NewServer(config ResourceHandlersConfig) *Server {
	router := NewRouter(config.Logger)

	resourceHandlers := &handlers.ResourceHandlers{
		URLShortnerHandler: handlers.NewURLShortnerHandler(config.Logger),
	}

	router.RegisterRoutes(resourceHandlers)
	return &Server{
		router: router,
		logger: config.Logger,
	}
}

func (s *Server) Start(port string) error {
	addr := fmt.Sprintf(":%s", port)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: CORS(s.router),
		// for image upload need more time out
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	s.logger.Info("Starting server", map[string]any{
		"port": port,
	})

	// return s.httpServer.ListenAndServeTLS(certFile, keyFile)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server", nil)
	return s.httpServer.Shutdown(ctx)
}
