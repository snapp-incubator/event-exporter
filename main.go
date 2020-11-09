package main

import (
	"context"
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	klog "k8s.io/klog/v2"
	"k8s.io/kubectl/pkg/util/logs"
)

func getK8sClient() *kubernetes.Clientset {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}
	return clientset
}

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	_, cancel := context.WithCancel(context.Background())
	clientset := getK8sClient()
	factory := informers.NewSharedInformerFactory(clientset, time.Hour*24)
	controller := NewEventExporterController(factory)
	stop := make(chan struct{})
	defer close(stop)
	go metricserver(cancel)
	err := controller.Run(stop)
	if err != nil {
		klog.Fatal(err)
	}
	select {}
}
