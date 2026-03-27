package middleware

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

// APMTransport is a performance monitoring interceptor used to record API duration
type APMTransport struct {
	Next http.RoundTripper
}

func (t *APMTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	resp, err := t.Next.RoundTrip(req)

	duration := time.Since(start)

	// Data can be subsequently reported to a real APM monitoring system (e.g., Prometheus, Datadog)
	status := "ERROR"
	if resp != nil {
		status = strconv.Itoa(resp.StatusCode)
	}
	log.Printf("\n[APM Metrics] Endpoint=%s | Method=%s | Status=%s | Duration=%v\n",
		req.URL.Path, req.Method, status, duration)

	return resp, err
}

// APM returns a TransportMiddleware
func APM() TransportMiddleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return &APMTransport{Next: next}
	}
}
