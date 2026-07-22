package cmd

import (
	"bytes"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/sincerefly/easycmd/internal/testutil"
	"github.com/spf13/pflag"
	v "github.com/spf13/viper"
)

func withIPFactory(t *testing.T, fake *testutil.FakeHTTPClient, chooser testutil.FakeChooser) (*bytes.Buffer, func()) {
	t.Helper()

	out := new(bytes.Buffer)
	oldFactory := ipServiceFactory
	ipServiceFactory = func(addrs []string) *IpService {
		return NewIpServiceWithDeps(addrs, fake, chooser, out)
	}
	return out, func() { ipServiceFactory = oldFactory }
}

func resetIPCommandFlags(t *testing.T) {
	t.Helper()
	flags := ipCmd.Flags()
	_ = flags.Set("all", "false")
	_ = flags.Set("random", "false")
	_ = flags.Set("server", "")
	flags.VisitAll(func(f *pflag.Flag) {
		f.Changed = false
	})
}

func runIPCommand(t *testing.T, args ...string) (testutil.CommandOutput, error) {
	t.Helper()
	resetIPCommandFlags(t)
	t.Cleanup(func() { rootCmd.SetArgs(nil) })
	cmdArgs := append([]string{"ip"}, args...)
	return testutil.ExecuteCommand(rootCmd, cmdArgs...)
}

func TestIPCommandDefault(t *testing.T) {
	resetConfigForTest(t)
	v.Set("ip.address", []string{"http://a", "http://b"})

	fake := testutil.NewFakeHTTPClient(map[string]testutil.FakeResponse{
		"http://a": {Body: []byte("1.2.3.4"), StatusCode: http.StatusOK},
	})
	out, restore := withIPFactory(t, fake, testutil.FakeChooser{Fixed: "http://a"})
	defer restore()

	output, err := runIPCommand(t)
	if err != nil {
		t.Fatalf("execute ip: %v", err)
	}
	if fake.CallCount() != 1 {
		t.Fatalf("expected 1 call, got %d", fake.CallCount())
	}
	if !strings.Contains(out.String(), "1.2.3.4") && !strings.Contains(output.Out, "1.2.3.4") {
		t.Fatalf("unexpected output: buf=%q cmd=%q", out.String(), output.Out)
	}
}

func TestIPCommandAll(t *testing.T) {
	resetConfigForTest(t)
	v.Set("ip.address", []string{"http://a", "http://b"})

	fake := testutil.NewFakeHTTPClient(map[string]testutil.FakeResponse{
		"http://a": {Body: []byte("1.1.1.1"), StatusCode: http.StatusOK},
		"http://b": {Body: []byte("2.2.2.2"), StatusCode: http.StatusOK},
	})
	out, restore := withIPFactory(t, fake, testutil.FakeChooser{})
	defer restore()

	_, err := runIPCommand(t, "-a")
	if err != nil {
		t.Fatalf("execute ip -a: %v", err)
	}
	if fake.CallCount() != 2 {
		t.Fatalf("expected 2 calls, got %d", fake.CallCount())
	}
	if !strings.Contains(out.String(), "1.1.1.1") || !strings.Contains(out.String(), "2.2.2.2") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestIPCommandRandomFlag(t *testing.T) {
	resetConfigForTest(t)
	v.Set("ip.address", []string{"http://a", "http://b"})

	fake := testutil.NewFakeHTTPClient(map[string]testutil.FakeResponse{
		"http://b": {Body: []byte("2.2.2.2"), StatusCode: http.StatusOK},
	})
	out, restore := withIPFactory(t, fake, testutil.FakeChooser{Fixed: "http://b"})
	defer restore()

	_, err := runIPCommand(t, "-r")
	if err != nil {
		t.Fatalf("execute ip -r: %v", err)
	}
	calls := fake.CallsSnapshot()
	if fake.CallCount() != 1 || calls[0].URL != "http://b" {
		t.Fatalf("unexpected calls: %+v", calls)
	}
	if !strings.Contains(out.String(), "2.2.2.2") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestIPCommandServer(t *testing.T) {
	resetConfigForTest(t)
	v.Set("ip.address", []string{"http://a"})

	fake := testutil.NewFakeHTTPClient(map[string]testutil.FakeResponse{
		"http://custom": {Body: []byte("9.9.9.9"), StatusCode: http.StatusOK},
	})
	out, restore := withIPFactory(t, fake, testutil.FakeChooser{})
	defer restore()

	_, err := runIPCommand(t, "-s", "http://custom")
	if err != nil {
		t.Fatalf("execute ip -s: %v", err)
	}
	calls := fake.CallsSnapshot()
	if fake.CallCount() != 1 || calls[0].URL != "http://custom" {
		t.Fatalf("unexpected calls: %+v", calls)
	}
	if !strings.Contains(out.String(), "9.9.9.9") {
		t.Fatalf("unexpected output: %q", out.String())
	}
}

func TestIPCommandEmptyServices(t *testing.T) {
	resetConfigForTest(t)
	v.Set("ip.address", []string{})

	fake := testutil.NewFakeHTTPClient(nil)
	_, restore := withIPFactory(t, fake, testutil.FakeChooser{Fixed: "http://a"})
	defer restore()

	_, err := runIPCommand(t)
	if !errors.Is(err, ErrNoServices) {
		t.Fatalf("expected ErrNoServices, got %v", err)
	}
	if fake.CallCount() != 0 {
		t.Fatalf("expected no HTTP calls, got %d", fake.CallCount())
	}
}
