package handlers

import (
	"net/http"

	"github.com/abdulmanafc2001/url-shortner/pkg/logger"
)

type URLShortnerHandler struct {
	logger *logger.Logger
}

func NewURLShortnerHandler(logger *logger.Logger) *URLShortnerHandler {
	return &URLShortnerHandler{
		logger: logger,
	}
}

func (h *URLShortnerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.URLShortner(w, r)
	}
}

func (h *URLShortnerHandler) URLShortner(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Hitting URL Shortner API", nil)
}
