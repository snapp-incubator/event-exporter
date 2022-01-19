package main

import (
	"fmt"

	"k8s.io/client-go/informers"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/tools/cache"

	v1 "k8s.io/api/core/v1"
)

// EventExporterController exports the kubernetes events that are added
type EventExporterController struct {
	informerFactory informers.SharedInformerFactory
	eventInformer   coreinformers.EventInformer
}

// Run starts shared informers and waits for the shared informer cache to
// synchronize.
func (c *EventExporterController) Run(stopCh <-chan struct{}) error {
	// Starts all the shared informers that have been created by the factory so far.
	c.informerFactory.Start(stopCh)
	// wait for the initial synchronization of the local cache.
	if !cache.WaitForCacheSync(stopCh, c.eventInformer.Informer().HasSynced) {
		return fmt.Errorf("failed to sync")
	}
	return nil
}

func (c *EventExporterController) eventAdd(obj interface{}) {
	event := obj.(*v1.Event)
	IncSummaryEvent(event)
	switch event.Type {
	case "Normal":
		IncNormalEvent(event)
	case "Warning":
		IncWarningEvent(event)
	}
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
