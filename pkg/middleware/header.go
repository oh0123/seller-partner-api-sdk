package middleware

import (
	"net/http"
)

// HeaderModifierTransport is used to uniformly inject custom headers into all outgoing requests
type HeaderModifierTransport struct {
	Headers map[string]string
	Next    http.RoundTripper
}

func (t *HeaderModifierTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// The RoundTripper specification requires that the input req cannot be directly modified; it must be cloned
	clone := req.Clone(req.Context())

	// Inject unified headers
	for k, v := range t.Headers {
		if v != "" {
			clone.Header.Set(k, v)
		}
	}

	return t.Next.RoundTrip(clone)
}

// HeaderModifier returns a TransportMiddleware
func HeaderModifier(headers map[string]string) TransportMiddleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return &HeaderModifierTransport{
			Headers: headers,
			Next:    next,
		}
	}
}
