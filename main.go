package main

import (
	"os"

	"github.com/kr/pretty"
	"github.com/skos-ninja/prometheus-rancher-discovery/services/rancher"
)

var hostname = os.Getenv("RANCHER_HOST")

func main() {
	client := rancher.NewClient(&hostname, nil)

	containers, err := client.GetContainers(nil)
	if err != nil {
		panic(err)
	}

	pretty.Print(containers)
}
