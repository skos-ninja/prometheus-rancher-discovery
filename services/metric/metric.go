package metric

import (
	"fmt"

	"discovery/services/rancher"
)

// Metric is the format that prometheus is expecting json data to be in for <file_sd_config>
type Metric struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

const pathLabelKey = "__metrics_path__"

// ConvertServiceStackToMetrics awd
func ConvertServiceStackToMetrics(serviceStack *map[string]map[string][]rancher.Container) ([]Metric, error) {
	metrics := []Metric{}

	for stackName, stack := range *serviceStack {
		fmt.Printf("Mapping services in %s\n", stackName)
		for serviceName, services := range stack {
			metric := Metric{
				Targets: []string{},
				Labels: map[string]string{
					"env": stackName,
					"job": serviceName,
				},
			}

			fmt.Printf("Mapping service %s\n", serviceName)

			for _, service := range services {
				target := getTargetFromContainer(service)

				// Set metric path from container label
				pathLabel := service.GetLabelOrDefault("path", "/system/metrics")
				if _, ok := metric.Labels[pathLabelKey]; ok {
					if metric.Labels[pathLabelKey] != pathLabel {
						return nil, fmt.Errorf("Differing metric paths: %s", serviceName)
					}
				}
				metric.Labels[pathLabelKey] = service.GetLabelOrDefault("path", "/system/metrics")

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
