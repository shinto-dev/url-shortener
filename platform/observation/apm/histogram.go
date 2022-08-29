package apm

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func NewHistogram(name string) *prometheus.HistogramVec {
	return promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: name,
		//Help:    "The total number of processed events",
		Buckets: []float64{0.01, 0.025, 0.05, .1, .25, .5, 1, 2, 3, 4, 5, 10},
	}, []string{"component", "name"})
}
