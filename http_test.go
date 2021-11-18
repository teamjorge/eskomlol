package eskomlol

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
)

type mockHTTPClient struct {
	data string
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	res := http.Response{}
	buf := bytes.NewBuffer([]byte(m.data))
	closer := io.NopCloser(buf)
	res.Body = closer

	return &res, nil
}

func TestDefaultRequest(t *testing.T) {
	ctx := context.Background()
	req, err := defaultRequest(ctx, "/blah", nil)
	if err != nil {
		t.Errorf("unexpected error while creating request: %v", err)
		return
	}

	if req.URL.String() != baseURL+"/blah" {
		t.Errorf("expected url to be %s, got %s", baseURL+"/blah", req.URL.String())
	}

	if req.Header.Get("User-Agent") != userAgent {
		t.Errorf("expected User-Agent header to be %s, got %s", userAgent, req.Header.Get("User-Agent"))
	}
}

func TestDoRequest(t *testing.T) {
	client := mockHTTPClient{
		data: "boo",
	}

	data, err := doRequest(context.Background(), &client, "/blah", nil)
	if err != nil {
		t.Errorf("unexpected error while performing request: %v", err)
		return
	}

	if string(data) != "boo" {
		t.Errorf("expected response data to be %s, got %s", "boo", string(data))
	}
}

func TestDoRequestJSON(t *testing.T) {
	client := mockHTTPClient{
		data: `{"thing": "yes"}`,
	}

	resItem := struct {
		Thing string `json:"thing,omitempty"`
	}{}

	err := doRequestJSON(context.Background(), &client, "/blah", nil, &resItem)
	if err != nil {
		t.Errorf("unexpected error while performing request: %v", err)
		return
	}

	if resItem.Thing != "yes" {
		t.Errorf("expected resItem.Thing to be %s, got %s", "yes", resItem.Thing)
	}
}

func TestGetClient(t *testing.T) {
	expectedClient := &mockHTTPClient{}
	c := Client{httpClient: expectedClient}

	client := getClient(&c)

	if expectedClient != client {
		t.Errorf("expected client to be %v, got %v", expectedClient, client)
	}
}
