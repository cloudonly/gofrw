// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

package main

import (
	"github.com/luomu/gofrw/hive/cell"
	"github.com/luomu/gofrw/metrics/metric"
)

var exampleMetricsCell = cell.Metric(newExampleMetrics)

type exampleMetrics struct {
	ExampleCounter metric.Counter
}

func newExampleMetrics() exampleMetrics {
	return exampleMetrics{
		ExampleCounter: metric.NewCounter(metric.CounterOpts{
			ConfigName: "misc_total",
			Namespace:  "cilium",
			Subsystem:  "example",
			Name:       "misc_total",
			Help:       "Counts miscellaneous events",
		}),
	}
}
