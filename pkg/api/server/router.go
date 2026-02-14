package server

import (
	"net/http"

	"github.com/abdulmanafc2001/url-shortner/pkg/api/handlers"
	"github.com/abdulmanafc2001/url-shortner/pkg/logger"
)

type Router struct {
	mux    *http.ServeMux
	logger *logger.Logger
}

func NewRouter(logger *logger.Logger) *Router {
	return &Router{
		mux:    http.NewServeMux(),
		logger: logger,
	}
}

func (r *Router) registerResourceRoutes(path string, handler http.Handler) {
	r.mux.Handle(path+"/", handler)
	r.mux.Handle(path, handler)
}

// RegisterRoutes registers all API routes
func (r *Router) RegisterRoutes(handlers *handlers.ResourceHandlers) {
	if handlers == nil {
		r.logger.Error("Resource handlers is nil", nil)
		return
	}

	if handlers.URLShortnerHandler != nil {
		r.registerResourceRoutes("/api/v1/url-shortner", handlers.URLShortnerHandler)
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Log request
	r.logger.Info("Incoming request", map[string]any{
		"method": req.Method,
		"path":   req.URL.Path,
	})

	// Set common headers
	w.Header().Set("Content-Type", "application/json")

	// Serve request
	r.mux.ServeHTTP(w, req)
}
