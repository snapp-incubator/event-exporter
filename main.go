package main

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	bctx := context.Background()
	ctx, cancel := context.WithCancel(bctx)
	defer cancel()

	clientset := getK8sClient()
	factory := informers.NewSharedInformerFactory(clientset, time.Hour*24)
	controller := NewEventExporterController(factory)
	err := controller.Run(ctx.Done())
	if err != nil {
		klog.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	startServer(":8090", mux, cancel)
}
