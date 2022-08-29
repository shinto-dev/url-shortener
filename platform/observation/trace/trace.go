package trace

import (
	"context"
)

type logFieldsKeyType string

const traceFieldKey = logFieldsKeyType("trace-key")

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceFieldKey, traceID)
}

func GetTraceID(ctx context.Context) string {
	traceID, ok := ctx.Value(traceFieldKey).(string)
	if !ok {
		return ""
	}

	return traceID
}
