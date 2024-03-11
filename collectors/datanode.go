package collectors

import (
	"log"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type DatanodeSubcollector struct {
	metrics map[string]*prometheus.Desc
}

func NewDatanodeSubcollector() Subcollector {
	return &DatanodeSubcollector{
		metrics: map[string]*prometheus.Desc{
			"BytesWritten": prometheus.NewDesc(
				"bytes_written",
				"Total number of bytes written to DataNode",
                []string{"hostname"}, prometheus.Labels{"type": "DataNodeActivity", "service": "DataNode"}),
			"BytesRead": prometheus.NewDesc(
				"bytes_read",
				"Total number of bytes read from DataNode",
                []string{"hostname"}, prometheus.Labels{"type": "DataNodeActivity", "service": "DataNode"}),
			"HeartbeatsAvgTime": prometheus.NewDesc(
				"heartbeats_avg_time",
				"Average heartbeat time in milliseconds",
                []string{"hostname"}, prometheus.Labels{"type": "DataNodeActivity", "service": "DataNode"}),
		},
	}
}

func (c *DatanodeSubcollector) Handles(modelerType string) bool {
    return strings.HasPrefix(modelerType, "DataNodeActivity")
}

func (c *DatanodeSubcollector) Collect(bean map[string]interface{}, ch chan<- prometheus.Metric) {
	hostname := bean["tag.Hostname"].(string)
	for key, value := range bean {
		if metric, ok := c.metrics[key]; ok {
            value, ok := value.(float64)
            if !ok {
                log.Printf("Failed to convert value for key %s to float64: %v", key, value)
                continue
            }
			metric, err := prometheus.NewConstMetric(metric, prometheus.GaugeValue, value, hostname)
			if err != nil {
				log.Printf("Failed to create metric %s: %e", metric, err)
			}
			ch <- metric
		}
	}
}
