# prometheus-rancher-discovery

This tool is designed to sit alongside a prometheus instance to allow it to discover services that are able to have their metrics scraped based off of labels on the stack in Rancher.

It uses the [rancher metadata service](https://rancher.com/docs/rancher/v1.6/en/rancher-services/metadata-service/) and the [prometheus file_sd_config](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#file_sd_config) to do this.

## Usage

Currently the tool requires it to be run as a sidekick to a prometheus container to allow it to share the filesystem to place the dynamic file for prometheus to discover.

Running is as simple as launching the docker container which can be found [here](https://hub.docker.com/r/skos/prometheus-rancher-discovery).
You can set the following env variables to help you configure it to your needs:
* `DEFAULT_FILE` is the output path for the generated config
* `RANCHER_HOST` is the full hostname (including http prefix) of the rancher metadata service

To make a service appear in the output you simply have to add the `ninja.skos.prometheus.rancher.port` with the internal port the container is running on and the system will do the rest!

You can also define the `ninja.skos.prometheus.rancher.path` with a custom metrics path to scrape from. The default is set to be `/system/metrics`

## Example output
```json
[
    {
        "targets": [
            "service-example:80"
        ],
        "labels": {
            "env": "dev",
            "__metrics_path__": "/system/metrics",
            "job": "service-example"
        }
    }
]
```
