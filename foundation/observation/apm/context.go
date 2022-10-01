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

func StartSegment(ctx context.Context, name string) *prometheus.Timer {
	hist := FromContext(ctx)
	return prometheus.NewTimer(hist.WithLabelValues("service", name))
}

func StartExternalSegment(ctx context.Context, name string) *prometheus.Timer {
	hist := FromContext(ctx)
	return prometheus.NewTimer(hist.WithLabelValues("external", name))
}

func StartDataStoreSegment(ctx context.Context, name string) *prometheus.Timer {
	hist := FromContext(ctx)
	return prometheus.NewTimer(hist.WithLabelValues("db", name))
}
