package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"discovery/services/metric"
	"discovery/services/rancher"
)

var hostname = os.Getenv("RANCHER_HOST")
var fileName = os.Getenv("DEFAULT_FILE")

type serviceMap struct {
	containers []*rancher.Container
}

func main() {
	client := rancher.NewClient(&hostname, nil)

	fmt.Printf("Saving file to %s\n", fileName)

	for {
		fmt.Println("Running metric fetch")
		metrics, err := getMetrics(client)
		if err != nil {
			panic(err)
		}

		fmt.Println("Running mertic save")
		writeMetrics(metrics, fileName)

		fmt.Println("Waiting for next run")
		// Wait 30 seconds before running again
		time.Sleep(30 * time.Second)
	}
}

func getMetrics(client *rancher.Client) ([]metric.Metric, error) {
	containers, err := client.GetContainers(nil)
	if err != nil {
		return nil, err
	}

	stackMap := make(map[string]map[string][]rancher.Container)

	for _, container := range containers {
		if container.System {
			continue
		}

		port := container.GetLabel("port")
		if port == nil {
			continue
		}

		if _, ok := stackMap[container.StackName]; !ok {
			stackMap[container.StackName] = make(map[string][]rancher.Container)
		}

		stackMap[container.StackName][container.ServiceName] = append(stackMap[container.StackName][container.ServiceName], container)
	}

	metrics, err := metric.ConvertServiceStackToMetrics(&stackMap)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func writeMetrics(metrics []metric.Metric, fileName string) error {
	content, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileName, content, 0644)
	if err != nil {
		return err
	}

	return nil
}
