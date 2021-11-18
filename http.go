package eskomlol

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	baseURL   string = "http://loadshedding.eskom.co.za/LoadShedding"
	userAgent string = "Mozilla/5.0 (X11; Linux x86_64; rv:69.0) Gecko/20100101 Firefox/69.0"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func defaultRequest(ctx context.Context, endpoint string, body io.Reader) (*http.Request, error) {
	requestURL := baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, body)
	if err != nil {
		return req, err
	}

	req.Header.Add("User-Agent", userAgent)

	return req, nil
}

func doRequest(ctx context.Context, client HttpClient, endpoint string, body io.Reader) ([]byte, error) {
	req, err := defaultRequest(ctx, endpoint, body)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)

	return data, err
}

func doRequestJSON(ctx context.Context, client HttpClient, endpoint string, body io.Reader, out interface{}) error {
	data, err := doRequest(ctx, client, endpoint, body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("%v. response: %s", err, string(data))
	}

	return nil
}

func getClient(c *Client) HttpClient {
	h := c.httpClient
	if h == nil {
		h = &http.Client{Timeout: c.timeout}
	}

	return h
}
