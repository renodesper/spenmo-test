package ctx

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type CtxKey int

const (
	// CtxRequestID ...
	CtxRequestID CtxKey = iota
)

// GetRequestID ...
func GetRequestID(ctx context.Context) string {
	res, _ := ctx.Value(CtxRequestID).(string)
	return res
}

// SetRequestID ...
func SetRequestID(ctx context.Context, val string) context.Context {
	if val == "" {
		val = uuid.New().String()
	}
	return context.WithValue(ctx, CtxRequestID, val)
}

// ExtractRequestID ...
func ExtractRequestID(ctx context.Context, r *http.Request) context.Context {
	ctx = SetRequestID(ctx, r.Header.Get("X-Request-Id"))
	return ctx
}
