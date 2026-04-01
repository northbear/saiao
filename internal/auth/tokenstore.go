package auth

type TokenStore struct{}

func NewTokenStore() *TokenStore {
	return &TokenStore{}
}

func (s *TokenStore) GroupForToken(_ string) (string, bool) {
	_ = s
	return "", false
}
