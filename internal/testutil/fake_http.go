package testutil

import (
	"errors"
	"maps"
	"sync"
)

type HTTPCall struct {
	URL     string
	Headers map[string]string
}

type FakeResponse struct {
	Body       []byte
	StatusCode int
	Err        error
}

type FakeHTTPClient struct {
	mu        sync.Mutex
	Responses map[string]FakeResponse
	Calls     []HTTPCall
}

func NewFakeHTTPClient(responses map[string]FakeResponse) *FakeHTTPClient {
	if responses == nil {
		responses = map[string]FakeResponse{}
	}
	return &FakeHTTPClient{Responses: responses}
}

func (f *FakeHTTPClient) Get(url string, headers map[string]string) ([]byte, int, error) {
	headerCopy := maps.Clone(headers)

	f.mu.Lock()
	f.Calls = append(f.Calls, HTTPCall{
		URL:     url,
		Headers: headerCopy,
	})
	f.mu.Unlock()

	resp, ok := f.Responses[url]
	if !ok {
		return nil, 0, errors.New("unexpected url: " + url)
	}
	return resp.Body, resp.StatusCode, resp.Err
}

func (f *FakeHTTPClient) CallCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.Calls)
}

func (f *FakeHTTPClient) CallsSnapshot() []HTTPCall {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]HTTPCall(nil), f.Calls...)
}
