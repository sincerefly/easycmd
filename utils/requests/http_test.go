package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != "test-agent" {
			t.Errorf("unexpected User-Agent: %s", r.Header.Get("User-Agent"))
		}
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	data, statusCode, err := Get(server.URL, map[string]string{"User-Agent": "test-agent"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if statusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", statusCode)
	}
	if string(data) != "ok" {
		t.Fatalf("expected ok, got %s", data)
	}
}

func TestGetNonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("fail"))
	}))
	defer server.Close()

	data, statusCode, err := Get(server.URL, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if statusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", statusCode)
	}
	if string(data) != "fail" {
		t.Fatalf("expected fail, got %s", data)
	}
}

func TestGetInvalidURL(t *testing.T) {
	_, statusCode, err := Get("://bad-url", nil)
	if err == nil {
		t.Fatal("expected error for invalid url")
	}
	if statusCode != 0 {
		t.Fatalf("expected status 0, got %d", statusCode)
	}
}

func TestPostSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		if payload["name"] != "david" {
			t.Fatalf("unexpected name: %v", payload["name"])
		}
		_, _ = w.Write([]byte(`{"json":{"name":"david"}}`))
	}))
	defer server.Close()

	body := map[string]string{"name": "david"}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}

	data, statusCode, err := Post(server.URL, bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if statusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", statusCode)
	}
	if !strings.Contains(string(data), "david") {
		t.Fatalf("unexpected response: %s", data)
	}
}

func TestClientGetWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ctx-ok"))
	}))
	defer server.Close()

	client := NewClient(time.Second)
	ctx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()

	data, statusCode, err := client.GetWithContext(ctx, server.URL, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if statusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", statusCode)
	}
	if string(data) != "ctx-ok" {
		t.Fatalf("expected ctx-ok, got %s", data)
	}
}

func TestClientPostWithContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("posted"))
	}))
	defer server.Close()

	client := NewClient(time.Second)
	ctx := t.Context()

	data, statusCode, err := client.PostWithContext(ctx, server.URL, strings.NewReader("payload"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if statusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", statusCode)
	}
	if string(data) != "posted" {
		t.Fatalf("expected posted, got %s", data)
	}
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func TestDoBodyReadError(t *testing.T) {
	client := &Client{
		httpClient: &http.Client{
			Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(errorReader{}),
					Header:     make(http.Header),
				}, nil
			}),
		},
	}

	_, statusCode, err := client.Get("http://example.com", nil)
	if err == nil {
		t.Fatal("expected read error")
	}
	if statusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", statusCode)
	}
}
