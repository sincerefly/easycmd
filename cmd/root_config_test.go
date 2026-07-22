package cmd

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	v "github.com/spf13/viper"
)

func testdataPath(name string) string {
	return filepath.Join("..", "internal", "testutil", "testdata", name)
}

func TestLoadConfigValidFile(t *testing.T) {
	resetConfigForTest(t)
	cfgFile = testdataPath("config_valid.toml")

	if err := loadConfig(); err != nil {
		t.Fatalf("loadConfig: %v", err)
	}

	addrs := v.GetStringSlice("ip.address")
	if len(addrs) != 2 || addrs[0] != "http://a" || addrs[1] != "http://b" {
		t.Fatalf("unexpected addresses: %v", addrs)
	}
	if !strings.HasPrefix(cfgFile, "Using config file:") {
		t.Fatalf("unexpected cfgFile message: %q", cfgFile)
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	resetConfigForTest(t)
	cfgFile = filepath.Join(t.TempDir(), "missing.toml")

	if err := loadConfig(); err != nil {
		t.Fatalf("loadConfig: %v", err)
	}
	if cfgFile != "No config file used" {
		t.Fatalf("expected no config message, got %q", cfgFile)
	}
}

func TestLoadConfigInvalidFile(t *testing.T) {
	resetConfigForTest(t)
	cfgFile = testdataPath("config_invalid.toml")

	err := loadConfig()
	if err == nil {
		t.Fatal("expected parse error")
	}
	if _, ok := errors.AsType[v.ConfigParseError](err); !ok {
		t.Fatalf("expected ConfigParseError, got %T: %v", err, err)
	}
}
