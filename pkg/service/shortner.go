package service

import (
	"crypto/sha256"
	"encoding/base64"
	"net/url"
	"strings"

	"github.com/abdulmanafc2001/url-shortner/pkg/storage"
)

type ShortenerService struct {
	store *storage.MemoryStore
}

func NewShortenerService(store *storage.MemoryStore) *ShortenerService {
	return &ShortenerService{store: store}
}

func (s *ShortenerService) Shorten(originalURL string) (string, error) {
	parsed, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return "", err
	}

	// if it's already exist , return that value from memory
	if code, ok := s.store.GetCode(originalURL); ok {
		return code, nil
	}

	code := generateShortCode(originalURL)
	domain := extractDomain(parsed.Host)

	s.store.Save(originalURL, code, domain)

	return code, nil
}

func (s *ShortenerService) Resolve(code string) (string, bool) {
	return s.store.GetURL(code)
}

func (s *ShortenerService) TopDomains(limit int) []DomainMetric {
	counts := s.store.GetDomainCounts()

	list := make([]DomainMetric, 0, len(counts))
	for domain, count := range counts {
		list = append(list, DomainMetric{
			Domain: domain,
			Count:  count,
		})
	}

	// Sort descending by count
	for i := 0; i < len(list); i++ {
		for j := i + 1; j < len(list); j++ {
			if list[j].Count > list[i].Count {
				list[i], list[j] = list[j], list[i]
			}
		}
	}

	if len(list) > limit {
		list = list[:limit]
	}

	return list
}

type DomainMetric struct {
	Domain string `json:"domain"`
	Count  int    `json:"count"`
}

func generateShortCode(input string) string {
	hash := sha256.Sum256([]byte(input))
	encoded := base64.URLEncoding.EncodeToString(hash[:])

	// remove special chars and keep it small
	code := strings.ReplaceAll(encoded, "-", "")
	code = strings.ReplaceAll(code, "_", "")

	return code[:8]
}

func extractDomain(host string) string {
	host = strings.ToLower(host)

	host = strings.TrimPrefix(host, "www.")

	return host
}
