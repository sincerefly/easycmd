package testutil

import (
	"net/http"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestFakeHTTPClient(t *testing.T) {
	fake := NewFakeHTTPClient(map[string]FakeResponse{
		"http://a": {Body: []byte("ok"), StatusCode: http.StatusOK},
	})

	body, code, err := fake.Get("http://a", map[string]string{"X-Test": "1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != http.StatusOK || string(body) != "ok" {
		t.Fatalf("unexpected response: code=%d body=%q", code, body)
	}
	if fake.CallCount() != 1 {
		t.Fatalf("expected 1 call, got %d", fake.CallCount())
	}

	calls := fake.CallsSnapshot()
	if calls[0].URL != "http://a" || calls[0].Headers["X-Test"] != "1" {
		t.Fatalf("unexpected call snapshot: %+v", calls[0])
	}
}

func TestFakeChooser(t *testing.T) {
	chooser := FakeChooser{Fixed: "http://fixed"}
	got, err := chooser.Choice([]string{"http://a", "http://b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "http://fixed" {
		t.Fatalf("expected fixed choice, got %q", got)
	}
}

func TestExecuteCommand(t *testing.T) {
	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("ran:%v", args)
		},
	}

	out, err := ExecuteCommand(cmd, "a", "b")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if !strings.Contains(out.Out, "ran:[a b]") {
		t.Fatalf("unexpected output: %q", out.Out)
	}
}

func TestSetCommandIO(t *testing.T) {
	cmd := &cobra.Command{}
	buf := new(strings.Builder)
	SetCommandIO(cmd, buf, buf)
	if cmd.OutOrStdout() != buf {
		t.Fatal("expected custom stdout writer")
	}
}
