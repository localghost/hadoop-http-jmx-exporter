package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jcmturner/gokrb5/spnego"
	"github.com/kostrzewa9ld/hadoop-http-jmx-exporter/collectors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/jcmturner/gokrb5.v7/client"
	"gopkg.in/jcmturner/gokrb5.v7/config"
	"gopkg.in/jcmturner/gokrb5.v7/keytab"
)

func main() {
	keytab, err := keytab.Load(os.Getenv("KERBEROS_KEYTAB_PATH"))
	if err != nil {
		log.Fatalf("Failed to load keytab: %e", err)
	}

	config, err := config.Load(os.Getenv("KERBEROS_CONFIG_PATH"))
	if err != nil {
		log.Fatalf("Failed to load config: %e", err)
	}

	urlsFromEnv := os.Getenv("JMX_URLS")
	if urlsFromEnv == "" {
		log.Fatalf("JMX_URLS env variable is not set")
	}
	urls := []url.URL{}
	for _, urlFromEnv := range strings.Split(urlsFromEnv, " ") {
		u, err := url.Parse(urlFromEnv)
		if err != nil {
			log.Fatalf("Failed to parse url: %e", err)
		}
		urls = append(urls, *u)
	}

	ktbClient := client.NewClientWithKeytab(os.Getenv("KERBEROS_PRINCIPAL"), os.Getenv("KERBEROS_REALM"), keytab, config)
    timeout := 60 * time.Second
    timeoutEnv := os.Getenv("HTTP_CLIENT_TIMEOUT_SECONDS")
    if timeoutEnv != "" {
        i, err := strconv.Atoi(timeoutEnv)
        if err != nil {
            log.Fatalf("Failed to parse HTTP_CLIENT_TIMEOUT_SECONDS: %e", err)
        }
        timeout = time.Duration(i) * time.Second
    }
	httpClient := http.Client{Timeout: timeout}

	collector, err := collectors.NewMetricsCollector(spnego.NewClient(ktbClient, &httpClient, ""), urls)
	if err != nil {
		log.Fatalf("Failed to create collector: %e", err)
	}
	registry := prometheus.NewRegistry()
	registry.Register(collector)

	address := os.Getenv("LISTEN_ADDRESS")
	if address == "" {
		address = "0.0.0.0"
	}
	port := os.Getenv("LISTEN_PORT")
	if port == "" {
		port = "9100"
	}
	log.Printf("Listening on %s:%s", address, port)
	http.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), nil))
}
