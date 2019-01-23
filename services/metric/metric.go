package metric

import (
	"fmt"

	"github.com/skos-ninja/prometheus-rancher-discovery/services/rancher"
)

// Metric is the format that prometheus is expecting json data to be in for <file_sd_config>
type Metric struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

// ConvertServiceStackToMetrics awd
func ConvertServiceStackToMetrics(serviceStack *map[string]map[string][]rancher.Container) ([]Metric, error) {
	metrics := []Metric{}

	for stackName, stack := range *serviceStack {
		fmt.Printf("Mapping services in %s\n", stackName)
		for serviceName, services := range stack {
			metric := Metric{
				Targets: []string{},
				Labels: map[string]string{
					"env":     stackName,
					"service": serviceName,
				},
			}

			fmt.Printf("Mapping service %s\n", serviceName)

			for _, service := range services {
				target := getTargetFromContainer(service)

				// Set metric path from container label
				pathLabel := service.GetLabelOrDefault("path", "/system/metrics")
				if _, ok := metric.Labels["metrics_path"]; ok {
					if metric.Labels["metrics_path"] != pathLabel {
						return nil, fmt.Errorf("Differing metric paths: %s", serviceName)
					}
				}
				metric.Labels["metrics_path"] = service.GetLabelOrDefault("path", "/system/metrics")

				metric.Targets = append(metric.Targets, target)
			}

			metrics = append(metrics, metric)
		}
	}

	return metrics, nil
}

func getTargetFromContainer(container rancher.Container) string {
	port := container.GetLabelOrDefault("port", "80")

	// We prefer the container name here rather than the
	return fmt.Sprintf("%s:%s", container.Name, port)
}
