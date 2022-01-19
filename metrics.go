package main

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	v1 "k8s.io/api/core/v1"
)

var (
	k8sNormalEvents = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "event",
		Subsystem: "normal",
		Name:      "k8s",
		Help:      "Exports kubernetes normal events count.",
	},
		[]string{
			"namespace",
			"reason",
			"kind",
			"source_host",
			"source_component",
		},
	)

	k8sSummaryEvents = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "event",
		Subsystem: "summary",
		Name:      "k8s",
		Help:      "Exports summary of k8s events count.",
	},
		[]string{
			"namespace",
			"reason",
			"kind",
			"type",
		},
	)

	k8sWarningEvents = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "event",
		Subsystem: "warning",
		Name:      "k8s",
		Help:      "Exports kubernetes warning events count.",
	},
		[]string{
			"namespace",
			"reason",
			"kind",
			"source_host",
			"source_component",
			"message",
		},
	)
)

func init() {
	prometheus.MustRegister(k8sNormalEvents)
	prometheus.MustRegister(k8sWarningEvents)
	prometheus.MustRegister(k8sSummaryEvents)
}

// IncSummaryEvent parses and increases an event counter with corresponding labels.
func IncSummaryEvent(event *v1.Event) {

	k8sSummaryEvents.WithLabelValues(
		event.Namespace,           // namespace
		event.Reason,              // reason
		event.InvolvedObject.Kind, // kind
		event.Type,                // type
	).Inc()
}

// IncNormalEvent parses and increases an event counter with corresponding labels.
func IncNormalEvent(event *v1.Event) {

	k8sNormalEvents.WithLabelValues(
		event.Namespace,           // namespace
		event.Reason,              // reason
		event.InvolvedObject.Kind, // kind
		event.Source.Host,         // source_host
		event.Source.Component,    // source_component
	).Inc()
}

// IncWarningEvent parses and increases an event counter with corresponding labels.
func IncWarningEvent(event *v1.Event) {
	m := ""

	if event.Reason == "FailedMount" {
		switch {
		case strings.Contains(event.Message, "timeout expired waiting for volumes to attach or mount"):
			m = "timeout expired waiting for volumes to attach or mount"
		case strings.Contains(event.Message, "rpc error: code = DeadlineExceeded desc = context deadline exceeded"):
			m = "rpc error: code = DeadlineExceeded desc = context deadline exceeded"
		case strings.Contains(event.Message, "volumeattachments.storage.k8s.io"):
			m = "volumeattachments not found"
		case strings.Contains(event.Message, ": secret") || strings.Contains(event.Message, ": configmap"):
			m = "secret or configmap error"
		}
	}

	if event.Reason == "FailedAttachVolume" {
		switch {
		case strings.Contains(event.Message, "is attached to a different instance"):
			m = "is attached to a different instance"
		case strings.Contains(event.Message, "Volume is already used by pod"):
			m = "Volume is already used by pod"
		case strings.Contains(event.Message, "Volume is already exclusively attached to one node"):
			m = "Volume is already exclusively attached to one node"
		case strings.Contains(event.Message, "attachment timeout for volume"):
			m = "attachment timeout"
		case strings.Contains(event.Message, "status must be available or downloading"):
			m = "status must be available or downloading"
		}
	}

	k8sWarningEvents.WithLabelValues(
		event.Namespace,           // namespace
		event.Reason,              // reason
		event.InvolvedObject.Kind, // kind
		event.Source.Host,         // source_host
		event.Source.Component,    // source_component
		m,                         // message
	).Inc()
}
