package requests

import (
	"context"
	"io"
	"net/http"
	"time"
)

const defaultTimeout = 30 * time.Second

// DefaultClient is the shared HTTP client used by package-level helpers.
var DefaultClient = NewClient(defaultTimeout)

type Client struct {
	httpClient *http.Client
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: timeout},
	}
}

func Get(url string, headers map[string]string) ([]byte, int, error) {
	return DefaultClient.Get(url, headers)
}

func Post(url string, body io.Reader) ([]byte, int, error) {
	return DefaultClient.Post(url, body)
}

func (c *Client) Get(url string, headers map[string]string) ([]byte, int, error) {
	return c.GetWithContext(context.Background(), url, headers)
}

func (c *Client) Post(url string, body io.Reader) ([]byte, int, error) {
	return c.PostWithContext(context.Background(), url, body)
}

func (c *Client) GetWithContext(ctx context.Context, url string, headers map[string]string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	for h, v := range headers {
		req.Header.Add(h, v)
	}
	return c.do(req)
}

func (c *Client) PostWithContext(ctx context.Context, url string, body io.Reader) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, 0, err
	}
	return c.do(req)
}

func (c *Client) do(req *http.Request) ([]byte, int, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return resBody, resp.StatusCode, nil
}
