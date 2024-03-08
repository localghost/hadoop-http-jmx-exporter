package collectors

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"sync"

	"github.com/jcmturner/gokrb5/spnego"
	"github.com/prometheus/client_golang/prometheus"
)

type Subcollector interface {
	Collect(bean map[string]interface{}, ch chan<- prometheus.Metric)
	Handles(modelerType string) bool
}

type MetricsCollector struct {
	client     *spnego.Client
	collectors []Subcollector
	urls       []url.URL
}

func NewMetricsCollector(client *spnego.Client, urls []url.URL) (*MetricsCollector, error) {
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
	resp, _ := c.client.Get(address.String())
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var asJson map[string]interface{}
	json.Unmarshal(body, &asJson)
	for _, bean := range asJson["beans"].([]interface{}) {
		bean := bean.(map[string]interface{})
		for _, collector := range c.collectors {
			if collector.Handles(bean["modelerType"].(string)) {
				collector.Collect(bean, ch)
			}
		}
	}
}
