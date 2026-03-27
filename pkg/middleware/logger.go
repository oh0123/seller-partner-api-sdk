package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"
)

var bufferPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

// Logger interface allows users to inject custom high-performance logging libraries (e.g., zap.SugaredLogger, logrus.Entry)
type Logger interface {
	Debugf(format string, v ...any)
	Infof(format string, v ...any)
	Errorf(format string, v ...any)
}

// defaultLogger is a simple implementation based on the standard library log
type defaultLogger struct {
	logger *log.Logger
}

func (l *defaultLogger) Debugf(format string, v ...any) {
	l.logger.Printf("[DEBUG] "+format, v...)
}

func (l *defaultLogger) Infof(format string, v ...any) {
	l.logger.Printf("[INFO] "+format, v...)
}

func (l *defaultLogger) Errorf(format string, v ...any) {
	l.logger.Printf("[ERROR] "+format, v...)
}

// Default output to Stdout Log
var DefaultLog = &defaultLogger{logger: log.New(os.Stdout, "", log.LstdFlags)}

// MaxLogBodyLength defines the maximum number of bytes/characters to log for HTTP request/response bodies.
// This can be modified by the caller to adjust the threshold for truncated logs.
var MaxLogBodyLength = 2048

// LoggerTransport is a high-performance logging interceptor.
// It eliminates the double memory allocation caused by io.ReadAll(Body),
// instead bypassing it to output logs when the underlying stream is actually read (Stream Tapping).
type LoggerTransport struct {
	Next   http.RoundTripper
	Logger Logger
}

func (t *LoggerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	logger := t.Logger
	if logger == nil {
		logger = DefaultLog
	}

	logger.Infof("--- [HTTP REQUEST] --- %s %s\n", req.Method, req.URL.String())

	// Use httputil.DumpRequestOut to print the full request (Headers + Body)
	dumpReq := req.Clone(req.Context())
	dumpBody := false
	if req.Body != nil && req.GetBody != nil {
		reqBody, _ := req.GetBody()
		// Limit the body read to prevent OOM
		dumpReq.Body = io.NopCloser(io.LimitReader(reqBody, int64(MaxLogBodyLength)))
		dumpBody = true
	}
	if dumpBytes, err := httputil.DumpRequestOut(dumpReq, dumpBody); err == nil {
		logger.Debugf("Request Details:\n%s\n", string(dumpBytes))
	}

	rsp, err := t.Next.RoundTrip(req)
	if err != nil {
		logger.Errorf("HTTP Request Failed: %v\n", err)
		return nil, err
	}

	if dumpBytes, err := httputil.DumpResponse(rsp, false); err == nil {
		logger.Debugf("--- [HTTP RESPONSE HEADERS] ---\n%s\n", string(dumpBytes))
	} else {
		logger.Infof("--- [HTTP RESPONSE] --- Status: %s\n", rsp.Status)
	}

	if rsp.Body != nil {
		buf := bufferPool.Get().(*bytes.Buffer)
		// Replace the real Response Body with a tapCloser that has Tee functionality.
		// When json.NewDecoder consumes streams, data is silently cloned into buf without allocating additional massive memory.
		rsp.Body = &tapCloser{
			r:      io.TeeReader(rsp.Body, buf),
			c:      rsp.Body,
			buf:    buf,
			logger: logger,
		}
	}

	return rsp, nil
}

// tapCloser implements io.ReadCloser
type tapCloser struct {
	r      io.Reader
	c      io.Closer
	buf    *bytes.Buffer
	logger Logger
}

func (t *tapCloser) Read(p []byte) (n int, err error) {
	return t.r.Read(p)
}

func (t *tapCloser) Close() error {
	// When consumption is complete at the bottom layer and Close() is called, print the massive Body as a Debug-level log
	printStr := t.buf.String()
	if len(printStr) > MaxLogBodyLength {
		printStr = printStr[:MaxLogBodyLength] + "...\n(truncated)"
	}
	t.logger.Debugf("Response Body: %s\n", printStr)
	
	// If using sync.Pool, you can reset it here: t.buf.Reset() to return to the object pool
	t.buf.Reset()
	bufferPool.Put(t.buf)

	return t.c.Close()
}

// Log returns a TransportMiddleware
func Log(logger Logger) TransportMiddleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return &LoggerTransport{
			Next:   next,
			Logger: logger,
		}
	}
}
