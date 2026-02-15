package utils

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/abdulmanafc2001/url-shortener/pkg/api/types"
)

func ValidateURLShorteningCreateReq(req *types.URLShortnerCreateReq) error {
	rawURL := strings.TrimSpace(req.URL)

	if rawURL == "" {
		return fmt.Errorf("url is required")
	}

	if len(rawURL) > 2048 {
		return fmt.Errorf("url too long")
	}

	if strings.Contains(rawURL, " ") {
		return fmt.Errorf("url must not contain spaces")
	}

	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("invalid url format")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("only http and https urls are allowed")
	}

	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("url must contain hostname")
	}

	// validate port
	if u.Port() != "" {
		port, err := strconv.Atoi(u.Port())
		if err != nil || port < 1 || port > 65535 {
			return fmt.Errorf("invalid port")
		}
	}

	// block localhost
	if host == "localhost" || host == "127.0.0.1" || host == "::1" {
		return fmt.Errorf("localhost urls are not allowed")
	}

	// block private IPs
	ip := net.ParseIP(host)
	if ip != nil {
		if ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() {
			return fmt.Errorf("private ip urls are not allowed")
		}
	}

	// basic domain check if not IP
	if ip == nil {
		if !strings.Contains(host, ".") {
			return fmt.Errorf("invalid hostname")
		}
	}

	return nil
}
