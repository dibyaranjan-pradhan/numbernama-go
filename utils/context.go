package utils

import (
	"context"
	"time"
)

// Timeout constants for external operations.
// Use these when creating context for DB, HTTP, or background tasks.
const (
	DefaultTimeout = 30 * time.Second
	DBTimeout      = 10 * time.Second
	HTTPTimeout    = 15 * time.Second
)

// WithTimeout creates a context with the given timeout.
// Caller must call the returned cancel function when done to release resources.
func WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeout)
}

// WithDBTimeout creates a context with the standard database operation timeout.
// Use for repository calls (FindOne, InsertOne, UpdateOne, etc.).
func WithDBTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return WithTimeout(ctx, DBTimeout)
}

// WithHTTPTimeout creates a context with the standard HTTP/external call timeout.
// Use for outbound HTTP, email, SMS, or third-party API calls.
func WithHTTPTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	return WithTimeout(ctx, HTTPTimeout)
}

// WithBackgroundTimeout creates a context from context.Background() with the given timeout.
// Use for fire-and-forget or startup operations that have no parent context
// (e.g. goroutines, constructors loading config).
func WithBackgroundTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// GetTraceID returns the trace ID from context (set by X-Trace-ID middleware or WithTraceID).
// Returns empty string if not set.
func GetTraceID(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	v := ctx.Value("traceID")
	if v == nil {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// WithTraceID returns a copy of ctx with the given trace ID.
// Used by middleware; handlers and services use GetTraceID(ctx) via logger.WithContext(ctx).
func WithTraceID(ctx context.Context, traceID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, "traceID", traceID)
}
