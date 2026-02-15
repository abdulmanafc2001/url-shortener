package tests

import (
	"testing"

	"github.com/abdulmanafc2001/url-shortener/pkg/service"
	"github.com/abdulmanafc2001/url-shortener/pkg/storage"
)

func TestSameURLReturnsSameCode(t *testing.T) {
	store := storage.NewMemoryStore()
	s := service.NewShortenerService(store)

	url := "https://youtube.com/watch?v=abc123"

	code1, err := s.Shorten(url)
	if err != nil {
		t.Fatal(err)
	}

	code2, err := s.Shorten(url)
	if err != nil {
		t.Fatal(err)
	}

	if code1 != code2 {
		t.Fatalf("expected same code but got %s and %s", code1, code2)
	}
}

func TestResolveURL(t *testing.T) {
	store := storage.NewMemoryStore()
	s := service.NewShortenerService(store)

	url := "https://stackoverflow.com/questions/123"

	code, err := s.Shorten(url)
	if err != nil {
		t.Fatal(err)
	}

	original, ok := s.Resolve(code)
	if !ok {
		t.Fatal("expected url to exist")
	}

	if original != url {
		t.Fatalf("expected %s but got %s", url, original)
	}
}

func TestTopDomains(t *testing.T) {
	store := storage.NewMemoryStore()
	s := service.NewShortenerService(store)

	s.Shorten("https://udemy.com/course/go")
	s.Shorten("https://udemy.com/course/k8s")
	s.Shorten("https://udemy.com/course/docker")
	s.Shorten("https://youtube.com/watch?v=1")
	s.Shorten("https://youtube.com/watch?v=2")
	s.Shorten("https://wikipedia.org/wiki/kubernetes")

	top := s.TopDomains(3)

	if len(top) != 3 {
		t.Fatalf("expected 3 results got %d", len(top))
	}

	if top[0].Domain != "udemy.com" {
		t.Fatalf("expected udemy.com first got %s", top[0].Domain)
	}
}
