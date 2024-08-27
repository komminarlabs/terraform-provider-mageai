package mageai

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

type Client interface {
	PipelineAPI() PipelineAPI
	Close()
}

type client struct {
	config ClientConfig
	apiURL *url.URL
}

func New(config *ClientConfig) (Client, error) {
	var err error
	c := &client{config: *config}

	hostAddress := config.Host
	if !strings.HasSuffix(config.Host, "/") {
		hostAddress = config.Host + "/"
	}

	c.apiURL, err = url.Parse(hostAddress)
	if err != nil {
		return nil, fmt.Errorf("parsing host URL: %w", err)
	}

	c.apiURL.Path = path.Join(c.apiURL.Path, "/api") + "/"

	if c.config.HTTPClient == nil {
		c.config.HTTPClient = &http.Client{Timeout: 10 * time.Second}
	}
	return c, nil
}

func (c *client) Close() {
	c.config.HTTPClient.CloseIdleConnections()
}

func (c *client) PipelineAPI() PipelineAPI {
	return c
}

func (c *client) makeAPICall(httpMethod, path string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(httpMethod, c.apiURL.String()+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", c.config.ApiKey)
	resp, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	return respBody, nil
}
