package cmd

import (
	"bytes"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/sincerefly/easycmd/internal/testutil"
)

func newTestIPService(t *testing.T, addrs []string, fake *testutil.FakeHTTPClient, chooser testutil.FakeChooser) (*IpService, *bytes.Buffer) {
	t.Helper()

	out := new(bytes.Buffer)
	service := NewIpServiceWithDeps(addrs, fake, chooser, out)
	return service, out
}

func TestIpServiceQueryRandom(t *testing.T) {
	fake := testutil.NewFakeHTTPClient(map[string]testutil.FakeResponse{
		"http://a": {Body: []byte("1.1.1.1"), StatusCode: http.StatusOK},
	})
	service, out := newTestIPService(t, []string{"http://a", "http://b"}, fake, testutil.FakeChooser{Fixed: "http://a"})

	if err := service.QueryRandom(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.CallCount() != 1 {
		t.Fatalf("expected 1 call, got %d", fake.CallCount())
	}
	calls := fake.CallsSnapshot()
	if calls[0].URL != "http://a" {
		t.Fatalf("expected http://a, got %q", calls[0].URL)
	}
	if calls[0].Headers["User-Agent"] != "Curl/7.55.1" {
		t.Fatalf("unexpected User-Agent: %q", calls[0].Headers["User-Agent"])
	}
	if got := strings.TrimSpace(out.String()); !strings.Contains(got, "1.1.1.1") {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestIpServiceQueryAll(t *testing.T) {
	fake := testutil.NewFakeHTTPClient(map[string]testutil.FakeResponse{
		"http://a": {Body: []byte("1.1.1.1"), StatusCode: http.StatusOK},
		"http://b": {Body: []byte("2.2.2.2"), StatusCode: http.StatusOK},
	})
	service, out := newTestIPService(t, []string{"http://a", "http://b"}, fake, testutil.FakeChooser{})

	if err := service.QueryAll(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fake.CallCount() != 2 {
		t.Fatalf("expected 2 calls, got %d", fake.CallCount())
	}
	if !strings.Contains(out.String(), "1.1.1.1") || !strings.Contains(out.String(), "2.2.2.2") {
		t.Fatalf("expected both IPs in output, got %q", out.String())
	}
}

func TestIpServiceQueryRandomEmptyServices(t *testing.T) {
	fake := testutil.NewFakeHTTPClient(nil)
	service, _ := newTestIPService(t, nil, fake, testutil.FakeChooser{Fixed: "http://a"})

	if err := service.QueryRandom(); !errors.Is(err, ErrNoServices) {
		t.Fatalf("expected ErrNoServices, got %v", err)
	}
	if fake.CallCount() != 0 {
		t.Fatalf("expected no HTTP calls, got %d", fake.CallCount())
	}
}

func TestIpServiceQueryAllEmptyServices(t *testing.T) {
	fake := testutil.NewFakeHTTPClient(nil)
	service, _ := newTestIPService(t, []string{}, fake, testutil.FakeChooser{})

	if err := service.QueryAll(); !errors.Is(err, ErrNoServices) {
		t.Fatalf("expected ErrNoServices, got %v", err)
	}
}

func TestIpServiceQueryServerIp(t *testing.T) {
	fake := testutil.NewFakeHTTPClient(map[string]testutil.FakeResponse{
		"http://custom": {Body: []byte("9.9.9.9"), StatusCode: http.StatusOK},
	})
	service, out := newTestIPService(t, []string{"http://a"}, fake, testutil.FakeChooser{})

	service.QueryServerIp("http://custom")
	if fake.CallCount() != 1 || fake.CallsSnapshot()[0].URL != "http://custom" {
		t.Fatalf("unexpected calls: %+v", fake.CallsSnapshot())
	}
	if !strings.Contains(out.String(), "9.9.9.9") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestIpServiceRequest(t *testing.T) {
	tests := []struct {
		name     string
		response testutil.FakeResponse
		want     string
	}{
		{
			name:     "network error",
			response: testutil.FakeResponse{Err: errors.New("network down")},
			want:     "(error: network down)",
		},
		{
			name:     "non ok status",
			response: testutil.FakeResponse{Body: []byte("fail"), StatusCode: http.StatusServiceUnavailable},
			want:     "(error: status 503)",
		},
		{
			name:     "success",
			response: testutil.FakeResponse{Body: []byte("8.8.8.8"), StatusCode: http.StatusOK},
			want:     "8.8.8.8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := testutil.NewFakeHTTPClient(map[string]testutil.FakeResponse{
				"http://svc": tt.response,
			})
			service, _ := newTestIPService(t, []string{"http://svc"}, fake, testutil.FakeChooser{})

			got := service.request("http://svc")
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
