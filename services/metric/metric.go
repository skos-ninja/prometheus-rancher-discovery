package metric

// Metric is the format that prometheus is expecting json data to be in for <file_sd_config>
type Metric struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}
