package collectors

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

type ClusterMetricsSubcollector struct {
	metrics map[string]*prometheus.Desc
}

func NewClusterMetricsSubcollector() Subcollector {
	return &ClusterMetricsSubcollector{
		metrics: map[string]*prometheus.Desc{
			"NumLostNMs": prometheus.NewDesc(
				"num_lost_nms",
				"Current number of lost NodeManagers for not sending heartbeats.",
                []string{"hostname"}, prometheus.Labels{"type": "ClusterMetrics", "service": "ResourceManager"}),
			"NumUnhealthyNMs": prometheus.NewDesc(
				"num_unhealthy_nms",
				"Current number of unhealthy NodeManagers",
                []string{"hostname"}, prometheus.Labels{"type": "ClusterMetrics", "service": "ResourceManager"}),
			"NumRebootedNMs": prometheus.NewDesc(
				"num_rebooted_nms",
				"Current number of rebooted NodeManagers",
                []string{"hostname"}, prometheus.Labels{"type": "ClusterMetrics", "service": "ResourceManager"}),
		},
	}
}

func (c *ClusterMetricsSubcollector) Handles(modelerType string) bool {
    return modelerType == "ClusterMetrics"
}

func (c *ClusterMetricsSubcollector) Collect(bean map[string]interface{}, ch chan<- prometheus.Metric) {
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
