package ctxs

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"time"
)

var ContextKeys = []string{UserInfoKey, UserTokenKey, UserSetTokenKey, UserRoleKey, MetadataKey, UserAppCodeKey, UserTenantCodeKey}

func CopyContext(ctx context.Context) context.Context {
	newCtx := context.Background()
	newCtx = trace.ContextWithSpanContext(newCtx, trace.SpanContextFromContext(ctx))
	for _, k := range ContextKeys {
		if v := ctx.Value(k); v != nil {
			newCtx = context.WithValue(newCtx, k, v)
		}
	}
	newCtx = SetUserCtx(newCtx, GetUserCtx(ctx))
	return newCtx
}

func GetDeadLine(ctx context.Context, defaultDeadLine time.Time) time.Time {
	dead, ok := ctx.Deadline()
	if !ok {
		return defaultDeadLine
	}
	return dead
}
