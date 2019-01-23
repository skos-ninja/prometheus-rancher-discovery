package rancher

import "fmt"

// Container is the object that is returned by the rancher metadata service
type Container struct {
	ID   string `json:"uuid"`
	Name string `json:"name"`

	ServiceName    string `json:"service_name"`
	ServiceID      string `json:"service_uuid"`
	ServiceReplica int    `json:"service_index"`

	StackName string `json:"stack_name"`
	StackID   string `json:"stack_uuid"`

	System bool `json:"system"`

	PrimaryIP         string `json:"primary_ip"`
	PrimaryMACAddress string `json:"primary_mac_address"`

	Labels map[string]string `json:"labels"`
}

// GetLabel will return the string value of the label based off of the prefix plus label
func (c *Container) GetLabel(label string) *string {
	labelName := generateLabelName(label)

	if i, ok := c.Labels[labelName]; ok {
		return &i
	}

	return nil
}

func generateLabelName(label string) string {
	return fmt.Sprintf("ninja.skos.prometheus.rancher.%s", label)
}
