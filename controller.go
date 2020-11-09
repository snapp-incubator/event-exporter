package main

import (
	"fmt"

	"k8s.io/client-go/informers"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/tools/cache"
	klog "k8s.io/klog/v2"

	v1 "k8s.io/api/core/v1"
)

// EventExporterController exports the kubernetes events that are added
type EventExporterController struct {
	informerFactory informers.SharedInformerFactory
	eventInformer   coreinformers.EventInformer
}

// Run starts shared informers and waits for the shared informer cache to
// synchronize.
func (c *EventExporterController) Run(stopCh chan struct{}) error {
	// Starts all the shared informers that have been created by the factory so far.
	c.informerFactory.Start(stopCh)
	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(stopCh, c.eventInformer.Informer().HasSynced) {
		return fmt.Errorf("Failed to sync")
	}
	return nil
}

func (c *EventExporterController) eventAdd(obj interface{}) {
	event := obj.(*v1.Event)
	klog.Infof("Event: %s/%s", event.Namespace, event.Name)
	ocEvents.WithLabelValues(
		event.Namespace,           // event_namespace
		event.Reason,              // event_reason
		event.InvolvedObject.Kind, // event_kind
		event.Type,                // event_type
		event.Source.Host,         // event_source_host
		event.Source.Component,    // event_source_component
	).Inc()
}

// NewEventExporterController creates a EventExporterController
func NewEventExporterController(informerFactory informers.SharedInformerFactory) *EventExporterController {
	eventInformer := informerFactory.Core().V1().Events()

	c := &EventExporterController{
		informerFactory: informerFactory,
		eventInformer:   eventInformer,
	}
	eventInformer.Informer().AddEventHandler(
		// Your custom resource event handlers.
		cache.ResourceEventHandlerFuncs{
			// Called on creation
			AddFunc: c.eventAdd,
		},
	)
	return c
}
