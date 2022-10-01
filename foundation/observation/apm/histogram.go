package apm

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	Component = "component"
	Operation = "name"
)

var metrics = make(map[string]*prometheus.HistogramVec)

func NewHistogram(name string) *prometheus.HistogramVec {
	if _, ok := metrics[name]; !ok {
		metrics[name] = promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    name,
			Buckets: []float64{0.01, 0.025, 0.05, .1, .25, .5, 1, 2, 3, 4, 5, 10},
		}, []string{Component, Operation})
	}

	return metrics[name]
}
