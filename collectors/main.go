package collectors

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/url"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/localghost/hadoop-http-jmx-exporter/httpclient"
)

type Subcollector interface {
	Collect(bean map[string]interface{}, ch chan<- prometheus.Metric)
	Handles(modelerType string) bool
}

type MetricsCollector struct {
	client     httpclient.HttpClient
	collectors []Subcollector
	urls       []url.URL
}

func NewMetricsCollector(client httpclient.HttpClient, urls []url.URL) (*MetricsCollector, error) {
	for _, u := range urls {
		q := u.Query()
		q.Set("qry", "Hadoop:service=*,name=*")
		u.RawQuery = q.Encode()
	}
	return &MetricsCollector{
		client: client,
		collectors: []Subcollector{
			NewJvmMetricsSubcollector(),
			NewDatanodeSubcollector(),
			NewClusterMetricsSubcollector(),
            NewNodeManagerMetricsSubcollector(),
		},
		urls: urls,
	}, nil
}

func (c *MetricsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c *MetricsCollector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	for _, u := range c.urls {
        wg.Add(1)
		go func(u url.URL) {
            defer wg.Done()
			c.collectFromUrl(u, ch)
		}(u)
	}
    wg.Wait()
}

func (c *MetricsCollector) collectFromUrl(address url.URL, ch chan<- prometheus.Metric) {
	resp, err := c.client.Get(address.String())
    if err != nil {
        log.Printf("Error getting: %s", address.String(), err)
        return
    }

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body from %s: %s", address.String(), err)
        return
    }

	var asJson map[string]interface{}
    if err := json.Unmarshal(body, &asJson); err != nil {
        log.Printf("Error unmarshalling JSON: %s", err)
        return
    }

    beans, ok := asJson["beans"].([]interface{})
    if !ok {
        log.Println("Error castin beans to list of interfaces")
        return
    }

	for _, bean := range beans {
		bean, ok := bean.(map[string]interface{})
        if !ok {
            log.Println("Error casting bean to map from string to interface")
            continue
        }
        modelerType, found := bean["modelerType"]
        if !found {
            log.Println("Did not find modelerType in bean. Skipping.")
            continue
        }
		for _, collector := range c.collectors {
			if collector.Handles(modelerType.(string)) {
				collector.Collect(bean, ch)
			}
		}
	}
}
