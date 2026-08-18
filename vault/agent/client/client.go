package client

import (
	"net/http"
	"sync"
)

type tokenRoundTripper struct {
	transport http.RoundTripper
	token     string
	mu        sync.RWMutex
}

func NewTokenRoundTripper(transport http.RoundTripper, token string) http.RoundTripper {
	return &tokenRoundTripper{
		transport: transport,
		token:     token,
	}
}

func (rt *tokenRoundTripper) SetToken(token string) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	rt.token = token
}

func (rt *tokenRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.mu.RLock()
	token := rt.token
	rt.mu.RUnlock()

	if token != "" {
		req.Header.Set("X-Vault-Token", token)
	}

	return rt.transport.RoundTrip(req)
}
