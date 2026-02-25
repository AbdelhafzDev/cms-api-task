package httpclient

import (
	"net/http"
	"time"
)

type Client struct {
	http    *http.Client
	baseURL string
	headers http.Header
}

type Config struct {
	BaseURL string
	Timeout time.Duration
	Headers map[string]string
}


func New(cfg *Config) *Client {
	timeout := 30 * time.Second
	if cfg != nil && cfg.Timeout > 0 {
		timeout = cfg.Timeout
	}

	c := &Client{
		http: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		headers: make(http.Header),
	}

	if cfg != nil {
		c.baseURL = cfg.BaseURL
		for k, v := range cfg.Headers {
			c.headers.Set(k, v)
		}
	}

	return c
}

func (c *Client) R(method, path string) *Request {
	return &Request{
		client:  c,
		method:  method,
		path:    path,
		headers: make(http.Header),
		query:   make(map[string]string),
	}
}
