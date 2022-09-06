package apm

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

type metricsFieldsKey string

const metricsFieldsKeyValue = metricsFieldsKey("key-2")

func WithAPM(ctx context.Context, metricsName string) context.Context {
	return context.WithValue(ctx, metricsFieldsKeyValue, NewHistogram(metricsName))
}

func FromContext(ctx context.Context) *prometheus.HistogramVec {
	return ctx.Value(metricsFieldsKeyValue).(*prometheus.HistogramVec)
}
