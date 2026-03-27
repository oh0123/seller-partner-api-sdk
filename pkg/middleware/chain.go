package middleware

import "net/http"

// TransportMiddleware defines a unified interceptor wrapper signature
type TransportMiddleware func(http.RoundTripper) http.RoundTripper

// Chain is used to connect multiple interceptors together like an onion model
func Chain(base http.RoundTripper, mws ...TransportMiddleware) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	// Iterate in reverse order to ensure the first passed middleware is wrapped on the outermost layer (executed first)
	for i := len(mws) - 1; i >= 0; i-- {
		base = mws[i](base)
	}
	return base
}
