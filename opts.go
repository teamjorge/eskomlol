package eskomlol

import (
	"time"
)

type ClientOpt func(*Client)

// WithTimeout sets the http timeout for the Client to the given duration.
func WithTimeout(timeout time.Duration) ClientOpt {
	return func(c *Client) {
		c.timeout = timeout
	}
}

func withHTTPClient(httpClient HttpClient) ClientOpt {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func withNowFunc(nowFunc func() time.Time) ClientOpt {
	return func(c *Client) {
		c.nowFunc = nowFunc
	}
}
