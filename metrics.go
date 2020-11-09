package main

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ocEvents = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "snappcloud",
		Subsystem: "event",
		Name:      "openshift",
		Help:      "Event that happend since this application become available.",
	},
		[]string{"event_namespace", "event_reason", "event_kind", "event_type", "event_source_host", "event_source_component"},
	)
)

func metricserver(cancel context.CancelFunc) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	startServer("8090", mux, cancel)
}

func init() {
	prometheus.MustRegister(ocEvents)
}
