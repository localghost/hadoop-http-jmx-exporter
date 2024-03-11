package collectors

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)


type JvmMetricsSubcollector struct {
	metrics map[string]*prometheus.Desc
}

func NewJvmMetricsSubcollector() Subcollector {
	return &JvmMetricsSubcollector{
		metrics: map[string]*prometheus.Desc{
			"MemNonHeapUsedM": prometheus.NewDesc(
				"jvm_mem_non_heap_used_m",
				"JVM non-heap used memory",
				[]string{"hostname", "service"}, prometheus.Labels{"type": "JvmMetrics"}),
			"MemNonHeapCommittedM": prometheus.NewDesc(
				"jvm_mem_non_heap_committed_m",
				"JVM non-heap committed memory",
				[]string{"hostname", "service"}, prometheus.Labels{"type": "JvmMetrics"}),
			"MemNonHeapMaxM": prometheus.NewDesc(
				"jvm_mem_non_heap_max_m",
				"JVM non-heap max memory",
				[]string{"hostname", "service"}, prometheus.Labels{"type": "JvmMetrics"}),
			"MemHeapUsedM": prometheus.NewDesc(
				"jvm_mem_heap_used_m",
				"JVM heap used memory",
				[]string{"hostname", "service"}, prometheus.Labels{"type": "JvmMetrics"}),
			"MemHeapCommittedM": prometheus.NewDesc(
				"jvm_mem_heap_committed_m",
				"JVM heap committed memory",
				[]string{"hostname", "service"}, prometheus.Labels{"type": "JvmMetrics"}),
			"MemHeapMaxM": prometheus.NewDesc(
				"jvm_mem_heap_max_m",
				"JVM heap max memory",
				[]string{"hostname", "service"}, prometheus.Labels{"type": "JvmMetrics"}),
			"MemMaxM": prometheus.NewDesc(
				"jvm_mem_max_m",
				"JVM max memory",
				[]string{"hostname", "service"}, prometheus.Labels{"type": "JvmMetrics"}),
			"GcCount": prometheus.NewDesc(
				"jvm_gc_count",
				"JVM garbage collection count",
				[]string{"hostname", "service"}, prometheus.Labels{"type": "JvmMetrics"}),
			"GcTimeMillis": prometheus.NewDesc(
				"jvm_gc_time_millis",
				"JVM garbage collection time in milliseconds",
				[]string{"hostname", "service"}, prometheus.Labels{"type": "JvmMetrics"}),
			"GcTotalExtraSleepTime": prometheus.NewDesc(
				"jvm_gc_total_extra_sleep_time",
				"JVM garbage collection total extra sleep time",
				[]string{"hostname", "service"}, prometheus.Labels{"type": "JvmMetrics"}),
		},
	}
}

func (c *JvmMetricsSubcollector) Handles(modelerType string) bool {
    return modelerType == "JvmMetrics"
}

func (c *JvmMetricsSubcollector) Collect(bean map[string]interface{}, ch chan<- prometheus.Metric) {
	processName := bean["tag.ProcessName"].(string)
	hostname := bean["tag.Hostname"].(string)
	for key, value := range bean {
		if metric, ok := c.metrics[key]; ok {
            value, ok := value.(float64)
            if !ok {
                log.Printf("Failed to convert value for key %s to float64: %v", key, value)
                continue
            }
			metric, err := prometheus.NewConstMetric(metric, prometheus.GaugeValue, value, hostname, processName)
			if err != nil {
				log.Printf("Failed to create metric %s: %e", metric, err)
			}
			ch <- metric
		}
	}
}
