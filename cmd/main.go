package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jcmturner/gokrb5/spnego"
	"github.com/kostrzewa9ld/hadoop-http-jmx-exporter/collectors"
	"github.com/kostrzewa9ld/hadoop-http-jmx-exporter/httpclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/jcmturner/gokrb5.v7/client"
	"gopkg.in/jcmturner/gokrb5.v7/config"
	"gopkg.in/jcmturner/gokrb5.v7/keytab"
)

type Config struct {
	HttpClientTimeoutSeconds time.Duration `yaml:"http_client_timeout_seconds" env:"HTTP_CLIENT_TIMEOUT_SECONDS" env-default:"60s"`
	KerberosPrincipal        string        `yaml:"kerberos_principal" env:"KERBEROS_PRINCIPAL"`
	KerberosRealm            string        `yaml:"kerberos_realm" env:"KERBEROS_REALM"`
	KerberosKeytabPath       string        `yaml:"kerberos_keytab_path" env:"KERBEROS_KEYTAB_PATH"`
	KerberosConfigPath       string        `yaml:"kerberos_config_path" env:"KERBEROS_CONFIG_PATH"`
	JmxUrls                  []string      `yaml:"jmx_urls" env:"JMX_URLS" env-required:""`
	ListenAddress            string        `yaml:"listen_address" env:"LISTEN_ADDRESS" env-default:"0.0.0.0"`
	ListenPort               int           `yaml:"listen_port" env:"LISTEN_PORT" env-default:"9100"`
}

func createHttpClient(cfg *Config) httpclient.HttpClient {
	noSpnegoClient := http.Client{Timeout: time.Duration(cfg.HttpClientTimeoutSeconds)}

	if cfg.KerberosPrincipal != "" {
		keytab, err := keytab.Load(cfg.KerberosKeytabPath)
		if err != nil {
			log.Fatalf("Failed to load keytab: %e", err)
		}

		config, err := config.Load(cfg.KerberosConfigPath)
		if err != nil {
			log.Fatalf("Failed to load config: %e", err)
		}

		ktbClient := client.NewClientWithKeytab(cfg.KerberosPrincipal, cfg.KerberosRealm, keytab, config)
		return httpclient.HttpClientWithSpnego{Client: spnego.NewClient(ktbClient, &noSpnegoClient, "")}
	}

	return httpclient.HttpClientPure{Client: &noSpnegoClient}
}

func readJMXUrls(cfg *Config) []url.URL {
	urls := []url.URL{}
	for _, urlFromEnv := range cfg.JmxUrls {
		u, err := url.Parse(urlFromEnv)
		if err != nil {
			log.Fatalf("Failed to parse url: %e", err)
		}
		urls = append(urls, *u)
	}
	return urls
}

func main() {
	var cfg Config
	if len(os.Args) > 1 {
		err := cleanenv.ReadConfig(os.Args[1], &cfg)
		if err != nil {
			log.Fatalf("Failed to read configuration file : %e", err)
		}
	} else {
		log.Println("config.yml not found reading from env variables only")
		err := cleanenv.ReadEnv(&cfg)
		if err != nil {
			log.Fatalf("Failed to read config: %e", err)
		}
	}

	collector, err := collectors.NewMetricsCollector(createHttpClient(&cfg), readJMXUrls(&cfg))
	if err != nil {
		log.Fatalf("Failed to create collector: %e", err)
	}
	registry := prometheus.NewRegistry()
	registry.Register(collector)

	log.Printf("Listening on %s:%s", cfg.ListenAddress, cfg.ListenPort)
	http.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.ListenAddress, cfg.ListenPort), nil))
}
