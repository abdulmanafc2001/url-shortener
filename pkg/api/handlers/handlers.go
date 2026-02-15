package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/abdulmanafc2001/url-shortner/pkg/api/types"
	"github.com/abdulmanafc2001/url-shortner/pkg/logger"
	"github.com/abdulmanafc2001/url-shortner/pkg/service"
	"github.com/abdulmanafc2001/url-shortner/utils"
)

type URLShortnerHandler struct {
	logger    *logger.Logger
	shortener *service.ShortenerService
	baseURL   string
}

func NewURLShortnerHandler(logger *logger.Logger, shortener *service.ShortenerService, baseURL string) *URLShortnerHandler {
	return &URLShortnerHandler{
		logger:    logger,
		shortener: shortener,
		baseURL:   baseURL,
	}
}

func (h *URLShortnerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateURLShort(w, r)
	}
}

func (h *URLShortnerHandler) CreateURLShort(w http.ResponseWriter, r *http.Request) {
	var req types.URLShortnerCreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request body", map[string]any{
			"error": err.Error(),
		})
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	h.logger.Info("Shortening URL", map[string]any{
		"url":    req.URL,
		"method": r.Method,
		"path":   r.URL.Path,
	})

	err := utils.ValidateURLShorteningCreateReq(&req)
	if err != nil {
		h.logger.Error("url shortening request validation failed", map[string]any{
			"error": err.Error(),
			"url":   req.URL,
		})

		utils.RespondWithError(w, http.StatusBadRequest, "url shortening request validation failed", err)
		return
	}

	code, err := h.shortener.Shorten(req.URL)
	if err != nil {
		h.logger.Error("invalid url", map[string]any{
			"error": err.Error(),
		})

		utils.RespondWithError(w, http.StatusBadRequest, "invalid url", err)
		return
	}

	resp := types.ShortenResponse{
		Code:     code,
		ShortURL: h.baseURL + "/" + code,
	}

	h.logger.Info("URL Shortened successfully", map[string]any{
		"code":     resp.Code,
		"shortURL": resp.ShortURL,
	})
	
	utils.RespondWithJSON(w, http.StatusOK, resp)
}
