package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	collector "sqs-exporter/collector"
	//collector "github.com/BruceLEO1969/sqs-exporter/collector"
)

func recordMetrics() {
	go func() {
		for {
			collector.MonitorSQS()
			time.Sleep(10 * time.Second)
		}
	}()
}

func main() {
	recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9434", nil)
}
