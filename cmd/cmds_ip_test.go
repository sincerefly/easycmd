package cmd

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

type mockHTTPClient struct {
	responses map[string]mockResponse
}

type mockResponse struct {
	body       []byte
	statusCode int
	err        error
}

func (m *mockHTTPClient) Get(url string, _ map[string]string) ([]byte, int, error) {
	resp, ok := m.responses[url]
	if !ok {
		return nil, 0, errors.New("unexpected url")
	}
	return resp.body, resp.statusCode, resp.err
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("copy stdout: %v", err)
	}
	_ = r.Close()
	return buf.String()
}

func TestIpServiceQueryRandom(t *testing.T) {
	client := &mockHTTPClient{
		responses: map[string]mockResponse{
			"http://a": {body: []byte("1.1.1.1"), statusCode: 200},
			"http://b": {body: []byte("2.2.2.2"), statusCode: 200},
		},
	}
	service := NewIpServiceWithClient([]string{"http://a", "http://b"}, client)

	out := captureStdout(t, func() {
		if err := service.QueryRandom(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(out, "1.1.1.1") && !strings.Contains(out, "2.2.2.2") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestIpServiceQueryAll(t *testing.T) {
	client := &mockHTTPClient{
		responses: map[string]mockResponse{
			"http://a": {body: []byte("1.1.1.1"), statusCode: 200},
			"http://b": {body: []byte("2.2.2.2"), statusCode: 200},
		},
	}
	service := NewIpServiceWithClient([]string{"http://a", "http://b"}, client)

	out := captureStdout(t, func() {
		if err := service.QueryAll(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	if !strings.Contains(out, "1.1.1.1") || !strings.Contains(out, "2.2.2.2") {
		t.Fatalf("expected both IPs in output, got %q", out)
	}
}

func TestIpServiceQueryAllConcurrentCompletion(t *testing.T) {
	client := &slowMockHTTPClient{}
	service := NewIpServiceWithClient([]string{"http://a", "http://b", "http://c"}, client)

	_ = captureStdout(t, func() {
		if err := service.QueryAll(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	if client.completed != 3 {
		t.Fatalf("expected 3 completed requests, got %d", client.completed)
	}
}

func TestIpServiceQueryRandomEmptyServices(t *testing.T) {
	service := NewIpServiceWithClient(nil, &mockHTTPClient{})
	if err := service.QueryRandom(); !errors.Is(err, ErrNoServices) {
		t.Fatalf("expected ErrNoServices, got %v", err)
	}
}

func TestIpServiceQueryAllEmptyServices(t *testing.T) {
	service := NewIpServiceWithClient([]string{}, &mockHTTPClient{})
	if err := service.QueryAll(); !errors.Is(err, ErrNoServices) {
		t.Fatalf("expected ErrNoServices, got %v", err)
	}
}

func TestIpServiceQueryServerIp(t *testing.T) {
	client := &mockHTTPClient{
		responses: map[string]mockResponse{
			"http://custom": {body: []byte("9.9.9.9"), statusCode: 200},
		},
	}
	service := NewIpServiceWithClient([]string{"http://a"}, client)

	out := captureStdout(t, func() {
		service.QueryServerIp("http://custom")
	})
	if !strings.Contains(out, "9.9.9.9") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestIpServiceRequestError(t *testing.T) {
	client := &mockHTTPClient{
		responses: map[string]mockResponse{
			"http://bad": {err: errors.New("network down")},
		},
	}
	service := NewIpServiceWithClient([]string{"http://bad"}, client)

	result := service.request("http://bad")
	if !strings.Contains(result, "network down") {
		t.Fatalf("expected error message in result, got %q", result)
	}
}

type slowMockHTTPClient struct {
	mu        sync.Mutex
	completed int
}

func (c *slowMockHTTPClient) Get(_ string, _ map[string]string) ([]byte, int, error) {
	c.mu.Lock()
	c.completed++
	c.mu.Unlock()
	return []byte("ip"), 200, nil
}
