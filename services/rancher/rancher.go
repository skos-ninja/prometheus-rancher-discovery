package rancher

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	pathlib "path"
	"time"
)

// Client is a rancher json client
type Client struct {
	Scheme string
	Host   string
	Prefix string

	UserAgent string

	Client *http.Client
}

var defaultUserAgent = "prometheus-rancher-discovery/0.1"
var defaultHostname = "http://rancher-metadata.rancher.internal/"

// NewClient creats a rancher client with a set of default params
func NewClient(baseURL *string, c *http.Client) *Client {
	if baseURL == nil {
		baseURL = &defaultHostname
	}

	remote, _ := url.Parse(*baseURL)

	if c == nil {
		c = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	return &Client{
		Scheme: remote.Scheme,
		Host:   remote.Host,
		Prefix: remote.Path,

		UserAgent: defaultUserAgent,

		Client: c,
	}
}

func (c *Client) createRequest(method string, path string) (*http.Request, error) {
	url := &url.URL{
		Scheme: c.Scheme,
		Host:   c.Host,
		Path:   pathlib.Join(c.Prefix, path),
	}

	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", defaultUserAgent)

	return req, nil
}

// GetContainers implements the 2015-12-19/containers method with minimal output properties
func (c *Client) GetContainers(ctx context.Context) ([]Container, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	req, err := c.createRequest("GET", "2015-12-19/containers")
	if err != nil {
		return nil, err
	}

	res, err := c.Client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, err
	}

	var containers []Container
	err = json.Unmarshal(body, &containers)
	if err != nil {
		return nil, err
	}

	return containers, nil
}
