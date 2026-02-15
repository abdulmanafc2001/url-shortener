package storage

import (
	"maps"
	"sync"
)

type MemoryStore struct {
	mu sync.RWMutex

	urlToCode map[string]string
	codeToUrl map[string]string

	domainCount map[string]int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		urlToCode:   make(map[string]string),
		codeToUrl:   make(map[string]string),
		domainCount: make(map[string]int),
	}
}

func (m *MemoryStore) GetCode(url string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	code, ok := m.urlToCode[url]
	return code, ok
}

func (m *MemoryStore) Save(url, code, domain string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.urlToCode[url] = code
	m.codeToUrl[code] = url
	m.domainCount[domain]++
}

func (m *MemoryStore) GetURL(code string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	url, ok := m.codeToUrl[code]
	return url, ok
}

func (m *MemoryStore) GetDomainCounts() map[string]int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	copyMap := make(map[string]int)

	maps.Copy(copyMap, m.domainCount)
	return copyMap
}
