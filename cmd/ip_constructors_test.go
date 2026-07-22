package cmd

import (
	"strings"
	"testing"

	"github.com/sincerefly/easycmd/internal/testutil"
)

func TestNewIpServiceConstructors(t *testing.T) {
	defaultService := NewIpService([]string{"http://a"})
	if defaultService == nil {
		t.Fatal("expected default service")
	}

	fake := testutil.NewFakeHTTPClient(nil)
	clientService := NewIpServiceWithClient([]string{"http://a"}, fake)
	if clientService == nil {
		t.Fatal("expected client service")
	}
}

func TestRandomChooser(t *testing.T) {
	chooser := randomChooser{}
	got, err := chooser.Choice([]string{"http://only"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "http://only" {
		t.Fatalf("expected http://only, got %q", got)
	}
}

func TestInitConfig(t *testing.T) {
	resetConfigForTest(t)
	cfgFile = testdataPath("config_valid.toml")

	initConfig()
	if !strings.HasPrefix(cfgFile, "Using config file:") {
		t.Fatalf("unexpected cfgFile message: %q", cfgFile)
	}
}

func TestQueryRandomChooserError(t *testing.T) {
	fake := testutil.NewFakeHTTPClient(nil)
	service := NewIpServiceWithDeps(
		[]string{"http://a"},
		fake,
		testutil.FakeChooser{Err: ErrNoServices},
		nil,
	)

	if err := service.QueryRandom(); err == nil {
		t.Fatal("expected chooser error")
	}
}
