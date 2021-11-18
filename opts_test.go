package eskomlol

import (
	"net/http"
	"testing"
	"time"
)

type fakeOptsHTTPClient struct{}

func (f *fakeOptsHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return nil, nil
}

func TestOpts(t *testing.T) {
	c := Client{}

	WithTimeout(5 * time.Second)(&c)

	if c.timeout != 5*time.Second {
		t.Errorf("expected timeout to be 5 seconds, got %v", c.timeout)
	}

	var calledFakeNow bool
	fakeNow := func() time.Time {
		calledFakeNow = true
		return time.Time{}
	}
	withNowFunc(fakeNow)(&c)

	c.nowFunc()

	if !calledFakeNow {
		t.Error("expected fakeNow to have been called")
	}

	httpClient := fakeOptsHTTPClient{}

	withHTTPClient(&httpClient)(&c)

	if c.httpClient != &httpClient {
		t.Errorf("expected httpClient to be %v, got: %v", &httpClient, c.httpClient)
	}
}
