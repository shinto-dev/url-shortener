package apm

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

type metricsFieldsKey string

const metricsFieldsKeyValue = metricsFieldsKey("key-2")

func WithValue(ctx context.Context, vec *prometheus.HistogramVec) context.Context {
	return context.WithValue(ctx, metricsFieldsKeyValue, vec)
}
