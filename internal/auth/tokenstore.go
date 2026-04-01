package auth

import "sync"

type TokenStore struct {
	mu     sync.RWMutex
	groups map[string]string
}

func NewTokenStore() *TokenStore {
	return &TokenStore{
		groups: make(map[string]string),
	}
}

func (s *TokenStore) Register(token, group string) {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.groups[token] = group
}

func (s *TokenStore) GroupForToken(token string) (string, bool) {
	if s == nil {
		return "", false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	group, ok := s.groups[token]
	return group, ok
}
