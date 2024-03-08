package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jcmturner/gokrb5/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	cpuTemp    *prometheus.GaugeVec
	hdFailures *prometheus.CounterVec
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		cpuTemp: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "cpu_temperature_celsius",
			Help: "Current temperature of the CPU.",
		}, []string{"node"}),
		hdFailures: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "hd_errors_total",
				Help: "Number of hard-disk errors.",
			},
			[]string{"device"},
		),
	}
	reg.MustRegister(m.cpuTemp)
	reg.MustRegister(m.hdFailures)
	return m
}

func createClient() *http.Client {
	ktbClient := client.NewClientWithKeytab(os.Getenv("KRB_USER"), os.Getenv("KRB_REALM"), )
	return &http.Client{}
}

func main() {
	// Create a non-global registry.
	reg := prometheus.NewRegistry()
    client := createClient()

	// Create new metrics and register them using the custom registry.
	m := NewMetrics(reg)
	// Set values for the new created metrics.
	m.cpuTemp.With(prometheus.Labels{"node": "one"}).Set(65.3)
	m.cpuTemp.With(prometheus.Labels{"node": "two"}).Set(-5.1)
	m.hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()

	// Expose metrics and custom registry via an HTTP server
	// using the HandleFor function. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(":9100", nil))
}
