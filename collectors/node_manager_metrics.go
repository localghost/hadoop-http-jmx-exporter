package collectors

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

type NodeManagerMetricsSubcollector struct {
	metrics map[string]*prometheus.Desc
}

func NewNodeManagerMetricsSubcollector() Subcollector {
	return &NodeManagerMetricsSubcollector{
		metrics: map[string]*prometheus.Desc{
			"ContainersLaunched": prometheus.NewDesc(
				"containers_launched",
				"Total number of launched containers",
                []string{"hostname"}, prometheus.Labels{"type": "NodeManagerMetrics", "service": "NodeManager"}),
			"ContainersCompleted": prometheus.NewDesc(
				"containers_completed",
				"Total number of successfully completed containers",
                []string{"hostname"}, prometheus.Labels{"type": "NodeManagerMetrics", "service": "NodeManager"}),
			"ContainersFailed": prometheus.NewDesc(
				"containers_failed",
				"Total number of failed containers",
                []string{"hostname"}, prometheus.Labels{"type": "NodeManagerMetrics", "service": "NodeManager"}),
			"ContainersKilled": prometheus.NewDesc(
				"containers_killed",
				"Total number of killed containers",
                []string{"hostname"}, prometheus.Labels{"type": "NodeManagerMetrics", "service": "NodeManager"}),
			"ContainersIniting": prometheus.NewDesc(
				"containers_initing",
				"Current number of initializing containers",
                []string{"hostname"}, prometheus.Labels{"type": "NodeManagerMetrics", "service": "NodeManager"}),
			"ContainersRunning": prometheus.NewDesc(
				"containers_running",
				"Current number of running containers",
                []string{"hostname"}, prometheus.Labels{"type": "NodeManagerMetrics", "service": "NodeManager"}),
			"AllocatedContainers": prometheus.NewDesc(
				"allocated_containers",
				"Current number of allocated containers",
                []string{"hostname"}, prometheus.Labels{"type": "NodeManagerMetrics", "service": "NodeManager"}),
			"AllocatedGB": prometheus.NewDesc(
				"allocated_gb",
				"Current allocated memory in GB",
                []string{"hostname"}, prometheus.Labels{"type": "NodeManagerMetrics", "service": "NodeManager"}),
			"AvailableGB": prometheus.NewDesc(
				"available_gb",
				"Current available memory in GB",
                []string{"hostname"}, prometheus.Labels{"type": "NodeManagerMetrics", "service": "NodeManager"}),
			"AllocatedVCores": prometheus.NewDesc(
				"allocated_vcores",
				"Current used vcores",
                []string{"hostname"}, prometheus.Labels{"type": "NodeManagerMetrics", "service": "NodeManager"}),
			"AvailableVCores": prometheus.NewDesc(
				"available_vcores",
				"Current available vcores",
                []string{"hostname"}, prometheus.Labels{"type": "NodeManagerMetrics", "service": "NodeManager"}),
			"ContainerLaunchDuration": prometheus.NewDesc(
				"container_launch_duration",
				"Average time duration in milliseconds NM takes to launch a container",
                []string{"hostname"}, prometheus.Labels{"type": "NodeManagerMetrics", "service": "NodeManager"}),
			"NodeUsedMemGB": prometheus.NewDesc(
				"node_used_mem_gb",
				"",
                []string{"hostname"}, prometheus.Labels{"type": "NodeManagerMetrics", "service": "NodeManager"}),
			"NodeUsedVMemGB": prometheus.NewDesc(
                "node_used_vmem_gb",
				"",
                []string{"hostname"}, prometheus.Labels{"type": "NodeManagerMetrics", "service": "NodeManager"}),
			"NodeCpuUtilization": prometheus.NewDesc(
                "node_cpu_utilization",
				"",
                []string{"hostname"}, prometheus.Labels{"type": "NodeManagerMetrics", "service": "NodeManager"}),
			"ContainerUsedMemGB": prometheus.NewDesc(
                "container_used_mem_gb",
				"",
                []string{"hostname"}, prometheus.Labels{"type": "NodeManagerMetrics", "service": "NodeManager"}),
			"ContainerUsedVMemGB": prometheus.NewDesc(
                "container_used_vmem_gb",
				"",
                []string{"hostname"}, prometheus.Labels{"type": "NodeManagerMetrics", "service": "NodeManager"}),
		},
	}
}

func (c *NodeManagerMetricsSubcollector) Handles(modelerType string) bool {
    return modelerType == "NodeManagerMetrics"
}

func (c *NodeManagerMetricsSubcollector) Collect(bean map[string]interface{}, ch chan<- prometheus.Metric) {
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
