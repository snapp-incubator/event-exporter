package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/build/kubernetes/api"
)

// Stream : structure for holding the stream of data coming from OpenShift
type Stream struct {
	Type  string    `json:"type,omitempty"`
	Event api.Event `json:"object"`
}

var (
	ocEvents = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "snappcloud",
		Subsystem: "event",
		Name:      "openshift",
		Help:      "Event that happend since this application become available.",
	},
		[]string{"event_namespace", "event_reason", "event_kind", "event_type", "event_message", "event_source_host", "event_source_component"},
	)
)

func metricserver() {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func init() {
	prometheus.MustRegister(ocEvents)
}

func main() {
	apiAddr := os.Getenv("OPENSHIFT_API_URL")
	apiToken := os.Getenv("OPENSHIFT_TOKEN")
	apiNamespace := os.Getenv("OPENSHIFT_NAMESPACE")
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c,
			syscall.SIGINT,  // Ctrl+C
			syscall.SIGTERM, // Termination Request
			syscall.SIGSEGV, // FullDerp
			syscall.SIGABRT, // Abnormal termination
			syscall.SIGILL,  // illegal instruction
			syscall.SIGFPE)  // floating point
		sig := <-c
		log.Fatalf("Signal (%v) Detected, Shutting Down", sig)
	}()
	go func() {
		metricserver()
	}()
	// check and make sure we have the minimum config information before continuing
	if apiAddr == "" {
		// use the default internal cluster URL if not defined
		apiAddr = "https://okd.private.teh-1.snappcloud.io"
		log.Print("Missing environment variable OPENSHIFT_API_URL. Using default API URL")
	}
	if apiToken == "" {
		// if we dont set it in the environment variable, read it out of
		// /var/run/secrets/kubernetes.io/serviceaccount/token
		log.Print("Missing environment variable OPENSHIFT_TOKEN. Leveraging serviceaccount token")
		fileData, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")

		if err != nil {
			log.Fatal("Service Account token does not exist.")
		}
		apiToken = string(fileData)
	}
	if apiNamespace == "" {
		log.Fatal("Missing environment variable OPENSHIFT_NAMESPACE. Exiting")
		os.Exit(1)
	}
	// setup ose connection
	var client http.Client
	client = http.Client{}
	req, err := http.NewRequest("GET", apiAddr+"/api/v1/namespaces/"+apiNamespace+"/events?watch=true", nil)
	if err != nil {
		log.Fatal("## Error while opening connection to openshift api", err)
	}
	req.Header.Add("Authorization", "Bearer "+apiToken)
	// req.Header.Add("Authorization", "Bearer cMGNrLGsPDf-_qsWE-4H3Y9zaAM0tmFsF98ISWNFdHw ")
	for {
		resp, err := client.Do(req)

		if err != nil {
			log.Println("## Error while connecting to:", apiAddr, err)
			time.Sleep(5 * time.Second)
			continue
		}

		streamStart := time.Now()
		reader := bufio.NewReader(resp.Body)

		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				log.Println("## Error reading from response stream.", err, line)
				resp.Body.Close()
				break
			}

			event := Stream{}
			decErr := json.Unmarshal(line, &event)
			if decErr != nil {
				log.Println("## Error decoding json.", err)
				resp.Body.Close()
				break
			}

			// Kubernetes sends all data from ETCD, we only want the logs since the stream started
			if event.Event.LastTimestamp.Time.After(streamStart) {
				ocEvents.WithLabelValues(
					event.Event.Namespace,           // event_namespace
					event.Event.Reason,              // event_reason
					event.Event.InvolvedObject.Kind, // event_kind
					event.Type,                      // event_type
					event.Event.Message,             // event_message
					event.Event.Source.Host,         // event_source_host
					event.Event.Source.Component,    // event_source_component
				).Inc()
			}

		}
	}
}
