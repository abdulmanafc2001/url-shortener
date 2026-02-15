package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/abdulmanafc2001/url-shortener/pkg/api/types"
	"github.com/abdulmanafc2001/url-shortener/pkg/logger"
	"github.com/abdulmanafc2001/url-shortener/pkg/service"
	"github.com/abdulmanafc2001/url-shortener/utils"
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
	trimmedURL := strings.Trim(r.URL.Path, "/")
	switch r.Method {
	case http.MethodPost:
		h.CreateURLShort(w, r)
	case http.MethodGet:
		if trimmedURL == "metrics" {
			h.Metrics(w, r)
			return
		}

		if trimmedURL == "" {
			utils.RespondWithError(w, http.StatusNotFound, "page not found", fmt.Errorf("page not found"))
			return
		}

		h.Redirect(w, r)

	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "method not allowed", fmt.Errorf("method not allowed"))
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
		ShortURL: h.baseURL + "?" + "code=" + code,
	}

	h.logger.Info("URL Shortened successfully", map[string]any{
		"code":     resp.Code,
		"shortURL": resp.ShortURL,
	})

	utils.RespondWithJSON(w, http.StatusOK, resp)
}

func (h *URLShortnerHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	original, ok := h.shortener.Resolve(code)
	if !ok {
		h.logger.Error("short url not found", nil)

		utils.RespondWithError(w, http.StatusNotFound, "short url not found", fmt.Errorf("short url not found"))
		return
	}

	http.Redirect(w, r, original, http.StatusFound)
}

func (h *URLShortnerHandler) Metrics(w http.ResponseWriter, r *http.Request) {
	top := h.shortener.TopDomains(3)

	utils.RespondWithJSON(w, http.StatusOK, top)
}
